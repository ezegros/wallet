package middlewares

import (
	"net/http"

	"github.com/ezegrosfeld/wallet/users/pkg/security"
	"github.com/gin-gonic/gin"
)

func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, nil)
			return
		}

		data, err := security.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, nil)
			return
		}

		c.Set("user", data)
	}
}
