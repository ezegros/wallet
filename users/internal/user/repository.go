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
	_, err = r.dynamo.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.table),
	})
	if err != nil {
		r.log.Errorw("Error inserting user", "error", err.Error())
		return err
	}

	return nil
}
