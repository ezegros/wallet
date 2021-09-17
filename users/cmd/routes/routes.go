package routes

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/ezegrosfeld/wallet/users/cmd/middlewares"
	"github.com/ezegrosfeld/wallet/users/internal/user"
	"github.com/ezegrosfeld/wallet/users/internal/wallet"
	"github.com/ezegrosfeld/wallet/users/pkg/db"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// MapRoutes creates all the routes needed
func MapRoutes(router *gin.Engine, dynamo dynamodbiface.DynamoDBAPI, log *zap.SugaredLogger) {
	// Create user repo, service and handler
	db.CreateTable(db.DynamoDB, "Users")
	ur := user.NewRepository(dynamo, "Users", log)
	us := user.NewService(log, ur)
	uh := user.NewHandler(us)

	// Create users router
	users := router.Group("/users")
	users.POST("/", uh.Create())
	users.POST("/login", uh.Login())

	// Create wallet repo, service and handler
	db.CreateTable(db.DynamoDB, "Wallets")
	wr := wallet.NewRepository(dynamo, "Wallets", log)
	ws := wallet.NewService(log, wr)
	wh := wallet.NewHandler(ws)

	// Create wallets router
	wallet := users.Group("/wallet")
	wallet.Use(middlewares.AuthorizationMiddleware())
	wallet.POST("/", wh.Create())
	wallet.GET("/address", wh.GetAddress())
}
