package routes

import (
	"testing"

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

// Test route mapping
func TestRouteMapping(t *testing.T) {
	router := gin.Default()

	db := &mockDynamoDBClient{}

	l, _ := zap.NewProduction()

	MapRoutes(router, db, l.Sugar())

	i := router.Routes()
	assert.Equal(t, len(i), 2)
}
