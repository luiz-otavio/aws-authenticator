package router

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	. "github.com/luiz-otavio/aws-authenticator/v2/internal"
	. "github.com/luiz-otavio/aws-authenticator/v2/pkg"
	"golang.org/x/crypto/bcrypt"
)

type changePasswordHandler struct {
	database Database
}

func NewChangePasswordHandler(database Database) ChangePasswordRouter {
	return changePasswordHandler{database: database}
}

func (impl changePasswordHandler) Handle(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var userRequest UserChangePasswordRequestSchema
	err := json.Unmarshal([]byte(request.Body), &userRequest)
	if err != nil {
		return nil, errors.New("failed to unmarshal request")
	}

	response, err := impl.ChangePassword(ctx, userRequest)
	if err != nil {
		return nil, errors.New("failed to change password")
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		return nil, errors.New("failed to marshal response")
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: response.StatusCode,
		Body:       string(responseBody),
	}, nil
}

func (impl changePasswordHandler) ChangePassword(ctx context.Context, request UserChangePasswordRequestSchema) (UserChangePasswordResponseSchema, error) {
	database := impl.database

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
				StatusCode: http.StatusInternalServerError,
				CommitedAt: time.Now(),
			},
			Message: "failed to get item",
		}, err
	}

	if output.Item == nil {
		return UserChangePasswordResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: http.StatusNotFound,
				CommitedAt: time.Now(),
			},
			Message: "user not found",
		}, nil
	}

	var targetUser User
	if err := dynamodbattribute.UnmarshalMap(output.Item, &targetUser); err != nil {
		return UserChangePasswordResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: http.StatusInternalServerError,
				CommitedAt: time.Now(),
			},
			Message: "failed to unmarshal item",
		}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.OldPassword), bcrypt.DefaultCost)
	if err != nil {
		return UserChangePasswordResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: http.StatusInternalServerError,
				CommitedAt: time.Now(),
			},
			Message: "failed to hash password",
		}, err
	}

	if !targetUser.Compare(hashedPassword) {
		return UserChangePasswordResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: http.StatusConflict,
				CommitedAt: time.Now(),
			},
			Message: "invalid password",
		}, errors.New("invalid password")
	}

	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return UserChangePasswordResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: http.StatusInternalServerError,
				CommitedAt: time.Now(),
			},
			Message: "failed to hash password",
		}, err
	}

	item, err := dynamodbattribute.MarshalMap(User{
		Username: request.Username,
		Password: hashedPassword,
	})

	if err != nil {
		return UserChangePasswordResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: http.StatusInternalServerError,
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
				StatusCode: http.StatusInternalServerError,
				CommitedAt: time.Now(),
			},
			Message: "failed to put item",
		}, err
	}

	return UserChangePasswordResponseSchema{
		HttpSchema: HttpSchema{
			StatusCode: http.StatusOK,
			CommitedAt: time.Now(),
		},
	}, nil
}
