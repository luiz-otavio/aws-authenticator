package cmd

import (
	"errors"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/credentials"
	. "github.com/luiz-otavio/aws-authenticator/v2/internal"
	"github.com/luiz-otavio/aws-authenticator/v2/internal/router"
	. "github.com/luiz-otavio/aws-authenticator/v2/pkg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Execute() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stderr,
	}).With().Timestamp().Logger()

	log.Debug().
		Time("Time", time.Now()).
		Msg("Initializing connect to DynamoDB...")

	region := os.Getenv("AWS_DYNAMODB_REGION")
	if len(region) == 0 {
		log.Error().
			Time("Time", time.Now()).
			Err(errors.New("Cannot find DynamoDB Region on env.")).
			Msg("Failed to connect to DynamoDB...")
		os.Exit(1)
		return
	}

	tableName := os.Getenv("AWS_DYNAMODB_TABLE_NAME")
	if len(tableName) == 0 {
		log.Error().
			Time("Time", time.Now()).
			Err(errors.New("Cannot find DynamoDB Table Name on env.")).
			Msg("Failed to connect to DynamoDB...")
		os.Exit(1)
		return
	}

	database, err := Init(credentials.NewCredentials(&credentials.EnvProvider{}), &region, &tableName)
	if err != nil {
		log.Error().
			Time("Time", time.Now()).
			Err(err).
			Msg("Failed to connect to DynamoDB...")
		os.Exit(1)
		return
	}

	log.Debug().
		Time("Time", time.Now()).
		Msg("Connected to DynamoDB successfully")

	log.Debug().
		Time("Time", time.Now()).
		Msg("Finding AWS_ROUTER...")

	awsRouter, err := GetRouter()
	if err != nil {
		log.Error().
			Time("Time", time.Now()).
			Err(err).
			Msg("Failed to find AWS_ROUTER...")
		os.Exit(1)
		return
	}

	log.Debug().
		Time("Time", time.Now()).
		Msg("AWS_ROUTER found successfully")

	var authHandler AuthHandler
	switch awsRouter {
	case LOGIN:
		authHandler = router.NewLoginHandler(database)
	case REGISTER:
		authHandler = router.NewRegisterHandler(database)
	case EXISTS:
		authHandler = router.NewExistsHandler(database)
	case CHANGE_PASSWORD:
		authHandler = router.NewChangePasswordHandler(database)
	default:
		log.Error().
			Time("Time", time.Now()).
			Err(errors.New("Invalid AWS_ROUTER")).
			Msg("Failed to find AWS_ROUTER...")
		os.Exit(1)
		return
	}

	log.Debug().
		Time("Time", time.Now()).
		Msg("AWS_ROUTER found successfully")

	lambda.Start(authHandler.Handle)
}
