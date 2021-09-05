package routes

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/ezegrosfeld/wallet/users/internal/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func MapRoutes(router *gin.Engine, dynamo dynamodbiface.DynamoDBAPI, log *zap.SugaredLogger) {
	// Create user repo
	ur := user.NewRepository(dynamo, "Users", log)
	us := user.NewService(log, ur)
	uh := user.NewHandler(us)

	users := router.Group("/users")
	users.POST("/", uh.Create())
}
