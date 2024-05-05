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

type registerHandler struct {
	database Database
}

func NewRegisterHandler(database Database) RegisterRouter {
	return registerHandler{database: database}
}

func (impl registerHandler) Handle(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var userRequest UserRegisterRequestSchema
	err := json.Unmarshal([]byte(request.Body), &userRequest)
	if err != nil {
		return nil, errors.New("failed to unmarshal request")
	}

	response, err := impl.Register(ctx, userRequest)
	if err != nil {
		return nil, errors.New("failed to register")
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

func (impl registerHandler) Register(ctx context.Context, request UserRegisterRequestSchema) (UserRegisterResponseSchema, error) {
	database := impl.database

	exists, err := database.DB().GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: database.TableName(),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: &request.Username,
			},
		},
	})

	if err != nil {
		return UserRegisterResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: 500,
				CommitedAt: time.Now(),
			},
			Message: "failed to get item",
		}, err
	}

	if exists.Item != nil {
		return UserRegisterResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: http.StatusConflict,
				CommitedAt: time.Now(),
			},
			Message: "username already exists",
		}, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserRegisterResponseSchema{
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
		return UserRegisterResponseSchema{
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
		return UserRegisterResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: http.StatusInternalServerError,
				CommitedAt: time.Now(),
			},
			Message: "failed to put item",
		}, err
	}

	return UserRegisterResponseSchema{
		HttpSchema: HttpSchema{
			StatusCode: http.StatusCreated,
			CommitedAt: time.Now(),
		},
		Message: "success",
	}, nil
}
