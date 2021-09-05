package db

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var DynamoDB *dynamodb.DynamoDB

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
		Credentials: credentials.NewSharedCredentials("/home/ezegrosfeld/.aws/credentials", "default"),
	})
	if err != nil {
		panic(err)
	}
	*/

	dynamodb := dynamodb.New(sess)

	DynamoDB = dynamodb
}
