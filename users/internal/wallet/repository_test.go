package wallet

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

func TestStore(t *testing.T) {
	mockClient := &mockDynamoDBClient{}
	mockClient.On("PutItemWithContext", mock.Anything, mock.Anything).Return(&dynamodb.PutItemOutput{}, nil)

	l, _ := zap.NewProduction()

	sl := l.Sugar()

	repo := NewRepository(mockClient, "Users", sl)

	wallet := &domain.Wallet{
		ID:     "123",
		UserID: "123",
		Seed:   "seed",
	}

	err := repo.Store(context.Background(), wallet)
	assert.Nil(t, err)
}

func TestStoreWithError(t *testing.T) {
	mockClient := &mockDynamoDBClient{}
	mockClient.On("PutItemWithContext", mock.Anything, mock.Anything).Return(&dynamodb.PutItemOutput{}, fmt.Errorf("error"))

	l, _ := zap.NewProduction()

	sl := l.Sugar()

	repo := NewRepository(mockClient, "Users", sl)

	wallet := &domain.Wallet{
		ID:     "123",
		UserID: "123",
		Seed:   "seed",
	}

	err := repo.Store(context.Background(), wallet)
	assert.Error(t, err)
}

func TestFind(t *testing.T) {
	mockClient := &mockDynamoDBClient{}

	o := &dynamodb.ScanOutput{
		Count: aws.Int64(1),
		Items: []map[string]*dynamodb.AttributeValue{
			{
				"id": {
					S: aws.String("123"),
				},
				"user_id": {
					S: aws.String("123"),
				},
				"seed": {
					S: aws.String("seed"),
				},
			},
		}}

	mockClient.On("Scan", mock.Anything).Return(o, nil)

	l, _ := zap.NewProduction()

	sl := l.Sugar()

	repo := NewRepository(mockClient, "Users", sl)

	wallet := &domain.Wallet{
		ID:     "123",
		UserID: "123",
		Seed:   "seed",
	}

	fw, err := repo.Find(context.Background(), wallet.UserID)
	assert.NoError(t, err)
	assert.Equal(t, wallet.ID, fw.ID)
	assert.Equal(t, wallet.UserID, fw.UserID)
	assert.Equal(t, wallet.Seed, fw.Seed)
}

func TestFindNotFound(t *testing.T) {
	mockClient := &mockDynamoDBClient{}

	o := &dynamodb.ScanOutput{
		Count: aws.Int64(1),
		Items: []map[string]*dynamodb.AttributeValue{}}

	mockClient.On("Scan", mock.Anything).Return(o, nil)

	l, _ := zap.NewProduction()

	sl := l.Sugar()

	repo := NewRepository(mockClient, "Users", sl)

	wallet := &domain.Wallet{
		ID:     "123",
		UserID: "123",
		Seed:   "seed",
	}

	_, err := repo.Find(context.Background(), wallet.UserID)
	assert.Error(t, err)
}

func TestFindWithError(t *testing.T) {
	mockClient := &mockDynamoDBClient{}

	o := &dynamodb.ScanOutput{
		Count: aws.Int64(1),
		Items: []map[string]*dynamodb.AttributeValue{
			{
				"id": {
					S: aws.String("123"),
				},
				"user_id": {
					S: aws.String("123"),
				},
				"seed": {
					S: aws.String("seed"),
				},
			},
		}}

	mockClient.On("Scan", mock.Anything).Return(o, fmt.Errorf("error"))

	l, _ := zap.NewProduction()

	sl := l.Sugar()

	repo := NewRepository(mockClient, "Users", sl)

	wallet := &domain.Wallet{
		ID:     "123",
		UserID: "123",
		Seed:   "seed",
	}

	_, err := repo.Find(context.Background(), wallet.UserID)
	assert.Error(t, err)
}
