package pkg

import "github.com/aws/aws-sdk-go/service/dynamodb"

type Database interface {
	TableName() *string
	DB() *dynamodb.DynamoDB
}
