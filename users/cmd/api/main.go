package main

import (
	"github.com/ezegrosfeld/wallet/users/cmd/routes"
	"github.com/ezegrosfeld/wallet/users/pkg/db"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Create a new gin router
	router := gin.Default()

	// Create a new logger
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	sl := l.Sugar()

	db.InitializeDatabase()

	dynamoDb := db.DynamoDB

	routes.MapRoutes(router, dynamoDb, sl)

	healthCheck(router)

	router.Run(":8080")
}

func healthCheck(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})
}
