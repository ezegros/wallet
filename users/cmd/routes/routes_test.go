package routes

/*package routes

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type mockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
	mock.Mock
}

func (m *mockDynamoDBClient) CreateTable(*dynamodb.CreateTableInput) (*dynamodb.CreateTableOutput, error) {
	args := m.Called()
	return args.Get(0).(*dynamodb.CreateTableOutput), args.Error(1)
}

// Test route mapping
func TestRouteMapping(t *testing.T) {
	router := gin.Default()

	db := &mockDynamoDBClient{}

	l, _ := zap.NewProduction()
	sl := l.Sugar()

	MapRoutes(router, db, sl)

	i := router.Routes()
	assert.Equal(t, len(i), 4)
}
*/
