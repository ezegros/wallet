package main

import (
	"github.com/ezegrosfeld/wallet/generator/cmd/routes"

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

	routes.MapRoutes(sl, router)

	healthCheck(router)

	router.Run(":8080")
}

func healthCheck(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})
}
