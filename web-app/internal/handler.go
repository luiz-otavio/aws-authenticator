package internal

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	. "github.com/luiz-otavio/aws-authenticator/v2/pkg"
)

type authHandler struct {
	AuthHandler
	AuthRouter
	Database
}

func Create(db Database) (authHandler, error) {
	return authHandler{}, nil
}

func (AuthHandler authHandler) Handle(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return nil, nil
}
