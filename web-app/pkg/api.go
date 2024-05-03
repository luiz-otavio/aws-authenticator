package pkg

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type AuthRouter interface {
	Login(context.Context, UserRequestLoginSchema) (UserResponseLoginSchema, error)
	Register(context.Context, UserRegisterRequestSchema) (UserRegisterResponseSchema, error)
	ChangePassword(context.Context, UserChangePasswordRequestSchema) (UserChangePasswordResponseSchema, error)
	Exists(context.Context, ExistsUserRequestSchema) (ExistsUserResponseSchema, error)
}

type AuthHandler interface {
	Handle(context.Context, *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)
}
