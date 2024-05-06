package router

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	. "github.com/luiz-otavio/aws-authenticator/v2/pkg"
)

type existsHandler struct {
	database Database
}

func NewExistsHandler(database Database) ExistsRouter {
	return existsHandler{database: database}
}

func (impl existsHandler) Handle(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var userRequest ExistsUserRequestSchema
	err := json.Unmarshal([]byte(request.Body), &userRequest)
	if err != nil {
		return nil, errors.New("failed to unmarshal request")
	}

	response, err := impl.Exists(ctx, userRequest)
	if err != nil {
		return nil, errors.New("failed to exists")
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

func (impl existsHandler) Exists(ctx context.Context, request ExistsUserRequestSchema) (ExistsUserResponseSchema, error) {
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
		return ExistsUserResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: http.StatusInternalServerError,
				CommitedAt: time.Now(),
			},
			Message: "failed to get item",
		}, err
	}

	if len(output.Item) == 0 {
		return ExistsUserResponseSchema{
			HttpSchema: HttpSchema{
				StatusCode: http.StatusNotFound,
				CommitedAt: time.Now(),
			},
			Message: "user not found",
		}, nil
	}

	return ExistsUserResponseSchema{
		HttpSchema: HttpSchema{
			StatusCode: http.StatusOK,
			CommitedAt: time.Now(),
		},
		Message: "user found",
	}, nil
}
