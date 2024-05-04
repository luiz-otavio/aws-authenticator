package internal

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	. "github.com/luiz-otavio/aws-authenticator/v2/pkg"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type authHandler struct {
	Database
}

func Create(db Database) (Authenticator, error) {
	if db.DB() == nil {
		return nil, errors.New("database is nil")
	}

	return authHandler{Database: db}, nil
}

func (handler authHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch event.HTTPMethod {
	case "POST":
		switch event.Path {
		case "/auth/login":
			var request UserRequestLoginSchema

			if err := json.Unmarshal([]byte(event.Body), &request); err != nil {
				log.Error().
					Err(err).
					Msg("failed to unmarshal request")

				return &events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       NewUnexpectedError(err, "failed to unmarshal request").String(),
				}, nil
			}

			response, err := handler.Login(ctx, request)
			if err != nil {
				log.Error().
					Err(err).
					Msg("failed to login")

				return &events.APIGatewayProxyResponse{
					StatusCode: response.StatusCode,
					Body:       NewUnexpectedError(err, "failed to login").String(),
				}, nil
			}

			return &events.APIGatewayProxyResponse{
				StatusCode: response.StatusCode,
				Body:       response.String(),
			}, nil
		case "/auth/register":
			var request UserRegisterRequestSchema
			if err := json.Unmarshal([]byte(event.Body), &request); err != nil {
				log.Error().
					Err(err).
					Msg("failed to unmarshal request")

				return &events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       NewUnexpectedError(err, "failed to unmarshal request").String(),
				}, nil
			}

			response, err := handler.Register(ctx, request)
			if err != nil {
				log.Error().
					Err(err).
					Msg("failed to register")

				return &events.APIGatewayProxyResponse{
					StatusCode: response.StatusCode,
					Body:       NewUnexpectedError(err, "failed to register").String(),
				}, nil
			}

			return &events.APIGatewayProxyResponse{
				StatusCode: response.StatusCode,
				Body:       response.String(),
			}, nil
		case "/auth/change-password":
			var request UserChangePasswordRequestSchema
			if err := json.Unmarshal([]byte(event.Body), &request); err != nil {
				log.Error().
					Err(err).
					Msg("failed to unmarshal request")

				return &events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       NewUnexpectedError(err, "failed to unmarshal request").String(),
				}, nil
			}

			response, err := handler.ChangePassword(ctx, request)
			if err != nil {
				log.Error().
					Err(err).
					Msg("failed to change password")

				return &events.APIGatewayProxyResponse{
					StatusCode: response.StatusCode,
					Body:       NewUnexpectedError(err, "failed to change password").String(),
				}, nil
			}

			return &events.APIGatewayProxyResponse{
				StatusCode: response.StatusCode,
				Body:       response.String(),
			}, nil
		default:
			return &events.APIGatewayProxyResponse{
				StatusCode: 404,
				Body:       NewUnexpectedError(nil, "invalid path").String(),
			}, nil
		}
	case "GET":
		switch event.Path {
		case "/auth/exists":
			var request ExistsUserRequestSchema
			if err := json.Unmarshal([]byte(event.Body), &request); err != nil {
				log.Error().
					Err(err).
					Msg("failed to unmarshal request")

				return &events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       NewUnexpectedError(err, "failed to unmarshal request").String(),
				}, nil
			}

			response, err := handler.Exists(ctx, request)
			if err != nil {
				log.Error().
					Err(err).
					Msg("failed to check if user exists")

				return &events.APIGatewayProxyResponse{
					StatusCode: response.StatusCode,
					Body:       NewUnexpectedError(err, "failed to check if user exists").String(),
				}, nil
			}

			return &events.APIGatewayProxyResponse{
				StatusCode: response.StatusCode,
				Body:       response.String(),
			}, nil
		default:
			return &events.APIGatewayProxyResponse{
				StatusCode: 404,
				Body:       NewUnexpectedError(nil, "invalid path").String(),
			}, nil
		}
	default:
		return &events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       NewUnexpectedError(nil, "invalid method").String(),
		}, nil
	}
}

func (authHandler authHandler) Login(ctx context.Context, request UserRequestLoginSchema) (UserResponseLoginSchema, error) {
	database := authHandler.Database

	output, err := database.DB().GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: database.TableName(),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: &request.Username,
			},
		},
	})

	if err != nil {
		return UserResponseLoginSchema{
			HttpSchema: HttpSchema{
				StatusCode: 500,
				CommitedAt: time.Now(),
			},
			Message: "failed to get item",
		}, err
	}

	var user user
	if err := dynamodbattribute.UnmarshalMap(output.Item, &user); err != nil {
		return UserResponseLoginSchema{
			HttpSchema: HttpSchema{
				StatusCode: 500,
				CommitedAt: time.Now(),
			},
			Message: "failed to unmarshal item",
		}, err
	}

	targetPassword := []byte(request.Password)
	if !user.Compare(targetPassword) {
		return UserResponseLoginSchema{
			HttpSchema: HttpSchema{
				StatusCode: 401,
				CommitedAt: time.Now(),
			},
			Message: "invalid password",
		}, errors.New("invalid password")
	}

	return UserResponseLoginSchema{
		HttpSchema: HttpSchema{
			StatusCode: 200,
			CommitedAt: time.Now(),
		},
		Message: "success",
	}, nil
}

func (authHandler authHandler) Register(ctx context.Context, request UserRegisterRequestSchema) (UserRegisterResponseSchema, error) {
	database := authHandler.Database

	exists, err := authHandler.Exists(ctx, ExistsUserRequestSchema{Username: request.Username})
	if err != nil {
		return UserRegisterResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: 500,
				CommitedAt: time.Now(),
			},
			Message: "failed to check if user exists",
		}, err
	}

	if exists.StatusCode == 200 {
		return UserRegisterResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: 409,
				CommitedAt: time.Now(),
			},
			Message: "user already exists",
		}, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserRegisterResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: 500,
				CommitedAt: time.Now(),
			},
			Message: "failed to hash password",
		}, err
	}

	item, err := dynamodbattribute.MarshalMap(user{
		username: request.Username,
		password: hashedPassword,
	})

	if err != nil {
		return UserRegisterResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: 500,
				CommitedAt: time.Now(),
			},
			Message: "failed to marshal item",
		}, err
	}

	if _, err := database.DB().PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: database.TableName(),
		Item:      item,
	}); err != nil {
		return UserRegisterResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: 500,
				CommitedAt: time.Now(),
			},
			Message: "failed to put item",
		}, err
	}

	return UserRegisterResponseSchema{
		HttpSchema: HttpSchema{
			StatusCode: 201,
			CommitedAt: time.Now(),
		},
		Message: "success",
	}, nil
}

func (handler authHandler) ChangePassword(ctx context.Context, request UserChangePasswordRequestSchema) (UserChangePasswordResponseSchema, error) {
	database := handler.Database

	if exists, err := handler.Exists(ctx, ExistsUserRequestSchema{Username: request.Username}); err != nil {
		return UserChangePasswordResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: 500,
				CommitedAt: time.Now(),
			},
			Message: "failed to check if user exists",
		}, err
	} else if exists.StatusCode != 200 {
		return UserChangePasswordResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: 404,
				CommitedAt: time.Now(),
			},
			Message: "user not found",
		}, errors.New("user not found")
	}

	var targetUser user
	output, err := database.DB().GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: database.TableName(),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: &request.Username,
			},
		},
	})

	if err != nil {
		return UserChangePasswordResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: 500,
				CommitedAt: time.Now(),
			},
			Message: "failed to get item",
		}, err
	}

	if err := dynamodbattribute.UnmarshalMap(output.Item, &targetUser); err != nil {
		return UserChangePasswordResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: 500,
				CommitedAt: time.Now(),
			},
			Message: "failed to unmarshal item",
		}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.OldPassword), bcrypt.DefaultCost)
	if err != nil {
		return UserChangePasswordResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: 500,
				CommitedAt: time.Now(),
			},
			Message: "failed to hash password",
		}, err
	}

	if !targetUser.Compare(hashedPassword) {
		return UserChangePasswordResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: 401,
				CommitedAt: time.Now(),
			},
			Message: "invalid password",
		}, errors.New("invalid password")
	}

	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return UserChangePasswordResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: 500,
				CommitedAt: time.Now(),
			},
			Message: "failed to hash password",
		}, err
	}

	item, err := dynamodbattribute.MarshalMap(user{
		username: request.Username,
		password: hashedPassword,
	})

	if err != nil {
		return UserChangePasswordResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: 500,
				CommitedAt: time.Now(),
			},
			Message: "failed to marshal item",
		}, err
	}

	if _, err := database.DB().PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: database.TableName(),
		Item:      item,
	}); err != nil {
		return UserChangePasswordResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: 500,
				CommitedAt: time.Now(),
			},
			Message: "failed to put item",
		}, err
	}

	return UserChangePasswordResponseSchema{
		HttpSchema: HttpSchema{
			StatusCode: 200,
			CommitedAt: time.Now(),
		},
	}, nil
}

func (handler authHandler) Exists(ctx context.Context, request ExistsUserRequestSchema) (ExistsUserResponseSchema, error) {
	database := handler.Database

	output, err := database.DB().GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: database.TableName(),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: &request.Username,
			},
		},
	})

	if err != nil {
		return ExistsUserResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: 500,
				CommitedAt: time.Now(),
			},
			Message: "failed to get item",
		}, err
	}

	if len(output.Item) == 0 {
		return ExistsUserResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: 404,
				CommitedAt: time.Now(),
			},
			Message: "user not found",
		}, nil
	}

	return ExistsUserResponseSchema{
		HttpSchema: HttpSchema{
			StatusCode: 200,
			CommitedAt: time.Now(),
		},
		Message: "user found",
	}, nil
}
