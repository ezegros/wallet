package wallet

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/ezegrosfeld/wallet/users/internal/domain"
	"go.uber.org/zap"
)

type Repository interface {
	Store(ctx context.Context, wallet *domain.Wallet) error
	Find(ctx context.Context, userID string) (*domain.Wallet, error)
}

type repository struct {
	// DynamoDB used for storage
	dynamo dynamodbiface.DynamoDBAPI
	// Name of the table
	table string
	// Logger
	log *zap.SugaredLogger
}

// NewRepository returns a new repository struct which must implement the Repository interface
func NewRepository(dynamo dynamodbiface.DynamoDBAPI, table string, log *zap.SugaredLogger) Repository {
	return &repository{
		dynamo: dynamo,
		table:  table,
		log:    log.Named("Wallet Repo"),
	}
}

// Store stores a wallet in the repository
func (r *repository) Store(ctx context.Context, wallet *domain.Wallet) error {
	av, err := dynamodbattribute.MarshalMap(wallet)
	if err != nil {
		r.log.Errorw("Error marshaling wallet", "error", err.Error())
		return err
	}

	// Put/Inser the item in the table
	_, err = r.dynamo.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.table),
	})
	if err != nil {
		r.log.Errorw("Error inserting wallet", "error", err.Error())
		return err
	}

	return nil
}

// Find finds a wallet in the repository by userID
func (r *repository) Find(ctx context.Context, userID string) (*domain.Wallet, error) {
	result, err := r.dynamo.Scan(&dynamodb.ScanInput{
		TableName:        aws.String(r.table),
		FilterExpression: aws.String("user_id = :user_id"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":user_id": {
				S: aws.String(userID),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, ErrNotFound
	}

	return itemToWallet(result.Items[0])

}

// Transforms the dynamo map into a usable domain.User
func itemToWallet(av map[string]*dynamodb.AttributeValue) (*domain.Wallet, error) {
	wallet := new(domain.Wallet)
	err := dynamodbattribute.UnmarshalMap(av, wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}
