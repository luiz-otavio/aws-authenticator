package pkg

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type LoginRouter interface {
	AuthHandler
	Login(context.Context, UserRequestLoginSchema) (UserResponseLoginSchema, error)
}

type RegisterRouter interface {
	AuthHandler
	Register(context.Context, UserRegisterRequestSchema) (UserRegisterResponseSchema, error)
}

type ChangePasswordRouter interface {
	AuthHandler
	ChangePassword(context.Context, UserChangePasswordRequestSchema) (UserChangePasswordResponseSchema, error)
}

type ExistsRouter interface {
	AuthHandler
	Exists(context.Context, ExistsUserRequestSchema) (ExistsUserResponseSchema, error)
}

type AuthHandler interface {
	Handle(context.Context, *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)
}

type Authenticator interface {
	AuthHandler
}
