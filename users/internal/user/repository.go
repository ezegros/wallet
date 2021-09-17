package user

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
	Store(ctx context.Context, user *domain.User) error
	Exists(username string) bool
	GetByUsername(username string) (*domain.User, error)
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
		log:    log.Named("User Repo"),
	}
}

// Store stores the user in the dynamodb table
func (r *repository) Store(ctx context.Context, user *domain.User) error {
	// Marshal the user struct into a map with dynamo attrs
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		r.log.Errorw("Error marshaling user", "error", err.Error())
		return err
	}

	// Put/Inser the item in the table
	_, err = r.dynamo.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.table),
	})
	if err != nil {
		r.log.Errorw("Error inserting user", "error", err.Error())
		return err
	}

	return nil
}

// Exists check if username us alredy registered
func (r *repository) Exists(username string) bool {
	// Get an item with username
	result, err := r.dynamo.Scan(&dynamodb.ScanInput{
		TableName:        aws.String(r.table),
		FilterExpression: aws.String("username = :username"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		return false
	}

	// Check if the amount of user is grater than 0
	return *result.Count >= 1
}

// GetByUsername returns the user found by username
func (r *repository) GetByUsername(username string) (*domain.User, error) {
	result, err := r.dynamo.Scan(&dynamodb.ScanInput{
		TableName:        aws.String(r.table),
		FilterExpression: aws.String("username = :username"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return itemToUser(result.Items[0])

}

// Transforms the dynamo map into a usable domain.User
func itemToUser(av map[string]*dynamodb.AttributeValue) (*domain.User, error) {
	user := new(domain.User)
	err := dynamodbattribute.UnmarshalMap(av, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
