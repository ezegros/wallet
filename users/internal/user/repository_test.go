package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
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

func (m *mockDynamoDBClient) PutItemWithContext(ctx context.Context, item *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error) {
	args := m.Called(item)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func (m *mockDynamoDBClient) GetItem(item *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	args := m.Called(item)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (m *mockDynamoDBClient) Scan(item *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	args := m.Called(item)
	return args.Get(0).(*dynamodb.ScanOutput), args.Error(1)
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

	db.On("PutItemWithContext", input).Return(o, nil)

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

	db.On("PutItemWithContext", input).Return(o, fmt.Errorf("something went wrong"))

	repo := NewRepository(db, "Users", sl)

	err := repo.Store(context.Background(), user)
	assert.Error(t, err)
}

func TestExists(t *testing.T) {
	db := &mockDynamoDBClient{}

	user := &domain.User{ID: "asda", Username: "username", Password: "contraseña"}

	input := &dynamodb.ScanInput{
		TableName:        aws.String("Users"),
		FilterExpression: aws.String("username = :username"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":username": {
				S: aws.String(user.Username),
			},
		},
	}

	l, _ := zap.NewProduction()

	sl := l.Sugar()

	o := &dynamodb.ScanOutput{
		Count: aws.Int64(1),
	}

	db.On("Scan", input).Return(o, nil)

	repo := NewRepository(db, "Users", sl)

	exists := repo.Exists(user.Username)
	assert.True(t, exists)
}

func TestNotExists(t *testing.T) {
	db := &mockDynamoDBClient{}

	user := &domain.User{ID: "asda", Username: "username", Password: "contraseña"}

	input := &dynamodb.ScanInput{
		TableName:        aws.String("Users"),
		FilterExpression: aws.String("username = :username"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":username": {
				S: aws.String(user.Username),
			},
		},
	}

	l, _ := zap.NewProduction()

	sl := l.Sugar()

	o := &dynamodb.ScanOutput{
		Count: aws.Int64(0),
	}

	db.On("Scan", input).Return(o, nil)

	repo := NewRepository(db, "Users", sl)

	exists := repo.Exists(user.Username)
	assert.False(t, exists)
}

func TestNotExistsWithError(t *testing.T) {
	db := &mockDynamoDBClient{}

	user := &domain.User{ID: "asda", Username: "username", Password: "contraseña"}

	input := &dynamodb.ScanInput{
		TableName:        aws.String("Users"),
		FilterExpression: aws.String("username = :username"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":username": {
				S: aws.String(user.Username),
			},
		},
	}

	l, _ := zap.NewProduction()

	sl := l.Sugar()

	o := &dynamodb.ScanOutput{
		Count: aws.Int64(0),
	}

	db.On("Scan", input).Return(o, fmt.Errorf("User not found"))

	repo := NewRepository(db, "Users", sl)

	exists := repo.Exists(user.Username)
	assert.False(t, exists)
}
