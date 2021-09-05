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
	dynamo dynamodbiface.DynamoDBAPI
	table  string
	log    *zap.SugaredLogger
}

func NewRepository(dynamo dynamodbiface.DynamoDBAPI, table string, log *zap.SugaredLogger) Repository {
	return &repository{
		dynamo: dynamo,
		table:  table,
		log:    log.Named("User Repo"),
	}
}

func (r *repository) Store(ctx context.Context, user *domain.User) error {
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		r.log.Errorw("Error marshaling user", "error", err.Error())
		return err
	}

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
