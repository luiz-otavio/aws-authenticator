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
)

type loginHandler struct {
	database Database
}

func NewLoginHandler(database Database) LoginRouter {
	return loginHandler{database: database}
}

func (impl loginHandler) Handle(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var userRequest UserRequestLoginSchema
	err := json.Unmarshal([]byte(request.Body), &userRequest)
	if err != nil {
		return nil, errors.New("failed to unmarshal request")
	}

	response, err := impl.Login(ctx, userRequest)
	if err != nil {
		return nil, errors.New("failed to login")
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

func (impl loginHandler) Login(ctx context.Context, request UserRequestLoginSchema) (UserResponseLoginSchema, error) {
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
		return UserResponseLoginSchema{
			HttpSchema: HttpSchema{
				StatusCode: http.StatusInternalServerError,
				CommitedAt: time.Now(),
			},
			Message: "failed to get item",
		}, err
	}

	var user User
	if err := dynamodbattribute.UnmarshalMap(output.Item, &user); err != nil {
		return UserResponseLoginSchema{
			HttpSchema: HttpSchema{
				StatusCode: http.StatusInternalServerError,
				CommitedAt: time.Now(),
			},
			Message: "failed to unmarshal item",
		}, err
	}

	targetPassword := []byte(request.Password)
	if !user.Compare(targetPassword) {
		return UserResponseLoginSchema{
			HttpSchema: HttpSchema{
				StatusCode: http.StatusUnauthorized,
				CommitedAt: time.Now(),
			},
			Message: "invalid password",
		}, errors.New("invalid password")
	}

	return UserResponseLoginSchema{
		HttpSchema: HttpSchema{
			StatusCode: http.StatusOK,
			CommitedAt: time.Now(),
		},
		Message: "success",
	}, nil
}
