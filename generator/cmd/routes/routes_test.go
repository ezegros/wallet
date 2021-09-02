package routes

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// Test route mapping
func TestRouteMapping(t *testing.T) {
	router := gin.Default()

	MapRoutes(&zap.SugaredLogger{}, router)

	i := router.Routes()
	assert.Equal(t, len(i), 2)
}
