package routes

import (
	address "github.com/ezegrosfeld/wallet/generator/internal/wallet"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func MapRoutes(log *zap.SugaredLogger, router *gin.Engine) {
	// Create address handler
	service := address.NewService(log)
	handler := address.NewHandler(service, log)

	// Map routes
	wallet := router.Group("/wallet")
	wallet.POST("/", handler.Create())
	wallet.GET("/", handler.Get())
}
