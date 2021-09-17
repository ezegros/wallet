package db

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var DynamoDB dynamodbiface.DynamoDBAPI

// InitializeDatabase intialize the dynamodb database with the aws SDK
func InitializeDatabase() {
	region := "us-west-2"
	endpoint := fmt.Sprintf("http://%s:%s", os.Getenv("DYNAMODB_SRV_SERVICE_HOST"), os.Getenv("DYNAMODB_SRV_SERVICE_PORT"))
	creds := credentials.NewStaticCredentials("local", "local", "")

	sess, err := session.NewSession(aws.NewConfig().WithCredentials(creds).WithRegion(region).WithEndpoint(endpoint))
	if err != nil {
		panic(err)
	}

	/* sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewSharedCredentials("/home/username/.aws/credentials", "default"),
	})
	if err != nil {
		panic(err)
	}
	*/

	dynamodb := dynamodb.New(sess)

	DynamoDB = dynamodb
}

func CreateTable(dynamo dynamodbiface.DynamoDBAPI, name string) error {
	_, err := dynamo.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(name),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})
	if err != nil {
		message := err.Error()
		if strings.Contains(message, "Cannot create preexisting table") {
			return nil
		}
		return err
	}
	return nil
}
