package internal

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	. "github.com/luiz-otavio/aws-authenticator/v2/pkg"
	"github.com/rs/zerolog/log"
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
					StatusCode: 500,
					Body:       NewUnexpectedError(err, "failed to login").String(),
				}, nil
			}

			return &events.APIGatewayProxyResponse{
				StatusCode: 200,
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
					StatusCode: 500,
					Body:       NewUnexpectedError(err, "failed to register").String(),
				}, nil
			}

			return &events.APIGatewayProxyResponse{
				StatusCode: 200,
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
					StatusCode: 500,
					Body:       NewUnexpectedError(err, "failed to change password").String(),
				}, nil
			}

			return &events.APIGatewayProxyResponse{
				StatusCode: 200,
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
					StatusCode: 500,
					Body:       NewUnexpectedError(err, "failed to check if user exists").String(),
				}, nil
			}

			return &events.APIGatewayProxyResponse{
				StatusCode: 200,
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
	return UserResponseLoginSchema{}, nil
}

func (authHandler authHandler) Register(ctx context.Context, request UserRegisterRequestSchema) (UserRegisterResponseSchema, error) {
	return UserRegisterResponseSchema{}, nil
}

func (handler authHandler) ChangePassword(ctx context.Context, request UserChangePasswordRequestSchema) (UserChangePasswordResponseSchema, error) {
	return UserChangePasswordResponseSchema{}, nil
}

func (handler authHandler) Exists(ctx context.Context, request ExistsUserRequestSchema) (ExistsUserResponseSchema, error) {
	return ExistsUserResponseSchema{}, nil
}
