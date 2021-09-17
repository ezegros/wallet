package db

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
	mock.Mock
}

func (m *mockDynamoDBClient) CreateTable(*dynamodb.CreateTableInput) (*dynamodb.CreateTableOutput, error) {
	args := m.Called()
	return args.Get(0).(*dynamodb.CreateTableOutput), args.Error(1)
}

func TestDatabaseInitialize(t *testing.T) {
	InitializeDatabase()

	db := DynamoDB

	assert.NotNil(t, db)
}

func TestCreateTable(t *testing.T) {
	mockClient := &mockDynamoDBClient{}
	mockClient.On("CreateTable", mock.Anything).Return(&dynamodb.CreateTableOutput{}, nil)

	DynamoDB = mockClient

	err := CreateTable(mockClient, "Users")

	assert.NoError(t, err)
}

func TestCreateTableError(t *testing.T) {
	mockClient := &mockDynamoDBClient{}
	mockClient.On("CreateTable", mock.Anything).Return(&dynamodb.CreateTableOutput{}, fmt.Errorf("error"))

	DynamoDB = mockClient

	err := CreateTable(mockClient, "Users")

	assert.Error(t, err)
}
