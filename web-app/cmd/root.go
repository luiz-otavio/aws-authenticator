package cmd

import (
	"errors"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/credentials"
	. "github.com/luiz-otavio/aws-authenticator/v2/internal"
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
	}

	log.Debug().
		Time("Time", time.Now()).
		Msg("Connected to DynamoDB successfully")

	handler, err := Create(database)
	if err != nil {
		log.Error().
			Time("Time", time.Now()).
			Err(err).
			Msg("Failed to initialize handler...")
		os.Exit(1)
		return
	}

	lambda.Start(handler.Handle)
}
