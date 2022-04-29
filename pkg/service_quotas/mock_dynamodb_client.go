package servicequotas

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockDynamoDbClient struct {
	dynamodbiface.DynamoDBAPI

	err               error
	ListTableResponse *dynamodb.ListTablesOutput
}
