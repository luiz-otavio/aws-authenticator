package pkg

import (
	"errors"
	"os"
)

type AWS_ROUTER string

const (
	LOGIN           AWS_ROUTER = "LOGIN"
	REGISTER        AWS_ROUTER = "REGISTER"
	EXISTS          AWS_ROUTER = "EXISTS"
	CHANGE_PASSWORD AWS_ROUTER = "CHANGE_PASSWORD"
)

func (e AWS_ROUTER) String() string {
	return string(e)
}

func GetRouter() (AWS_ROUTER, error) {
	router := os.Getenv("AWS_ROUTER")
	if router == "" {
		return "", errors.New("AWS_ROUTER not set")
	}

	switch router {
	case "LOGIN":
		return LOGIN, nil
	case "REGISTER":
		return REGISTER, nil
	case "EXISTS":
		return EXISTS, nil
	case "CHANGE_PASSWORD":
		return CHANGE_PASSWORD, nil
	default:
		return "", errors.New("invalid AWS_ROUTER")
	}
}
