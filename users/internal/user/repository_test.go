package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/ezegrosfeld/wallet/users/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type mockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
	mock.Mock
}

func (m *mockDynamoDBClient) PutItem(item *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	args := m.Called(item)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func TestCreate(t *testing.T) {
	db := &mockDynamoDBClient{}

	user := &domain.User{ID: "asda", Username: "username", Password: "contraseña"}

	av, _ := dynamodbattribute.MarshalMap(user)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Users"),
	}

	l, _ := zap.NewProduction()

	sl := l.Sugar()

	o := &dynamodb.PutItemOutput{}

	db.On("PutItem", input).Return(o, nil)

	repo := NewRepository(db, "Users", sl)

	err := repo.Store(context.Background(), user)
	assert.NoError(t, err)
}

func TestCreateWithFail(t *testing.T) {
	db := &mockDynamoDBClient{}

	user := &domain.User{ID: "asda", Username: "username", Password: "contraseña"}

	av, _ := dynamodbattribute.MarshalMap(user)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Users"),
	}

	l, _ := zap.NewProduction()

	sl := l.Sugar()

	o := &dynamodb.PutItemOutput{}

	db.On("PutItem", input).Return(o, fmt.Errorf("something went wrong"))

	repo := NewRepository(db, "Users", sl)

	err := repo.Store(context.Background(), user)
	assert.Error(t, err)
}
