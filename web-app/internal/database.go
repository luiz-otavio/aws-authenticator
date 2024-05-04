package internal

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	. "github.com/luiz-otavio/aws-authenticator/v2/pkg"
	"github.com/rs/zerolog/log"
)

type database struct {
	_tableName *string
	_db        *dynamodb.DynamoDB
}

func Init(credentials *credentials.Credentials, region *string, tableName *string) (Database, error) {
	session, err := session.NewSession(&aws.Config{
		Credentials: credentials,
		Region:      region,
	})

	if err != nil {
		return nil, err
	}

	var db = dynamodb.New(session)
	req, resp := db.DescribeTableRequest(&dynamodb.DescribeTableInput{TableName: tableName})

	err = req.Send()
	if err != nil {
		return nil, err
	}

	log.Debug().
		Time("Time", time.Now()).
		Msg(resp.String())

	log.Debug().
		Time("Time", time.Now()).
		Msg("Checking DynamoDB connection...")

	log.Debug().
		Time("Time", time.Now()).
		Msg("DynamoDB Service Name: " + db.ServiceName)

	return database{_tableName: tableName, _db: db}, nil
}

func (impl database) DB() *dynamodb.DynamoDB {
	return impl._db
}

func (impl database) TableName() *string {
	return impl._tableName
}
