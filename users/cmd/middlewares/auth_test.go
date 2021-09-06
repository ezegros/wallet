package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ezegrosfeld/wallet/users/pkg/security"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddlewareNoCookie(t *testing.T) {
	router := gin.Default()

	router.Use(AuthorizationMiddleware())

	router.GET("/authrequired", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, nil)
	})

	r := httptest.NewRequest("GET", "/authrequired", http.NoBody)
	rw := httptest.NewRecorder()

	router.ServeHTTP(rw, r)

	assert.Equal(t, http.StatusUnauthorized, rw.Code)
}

func TestAuthMiddlewareCookieWithInvalidJWT(t *testing.T) {
	router := gin.Default()

	router.Use(AuthorizationMiddleware())

	router.GET("/authrequired", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, nil)
	})

	r := httptest.NewRequest("GET", "/authrequired", http.NoBody)

	r.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "wdadljashflksdf",
	})

	rw := httptest.NewRecorder()

	router.ServeHTTP(rw, r)

	assert.Equal(t, http.StatusUnauthorized, rw.Code)
}

func TestAuthMiddleware(t *testing.T) {
	router := gin.Default()

	router.Use(AuthorizationMiddleware())

	router.GET("/authrequired", func(c *gin.Context) {
		data := c.MustGet("user").(*security.JWTData)
		c.JSON(http.StatusAccepted, gin.H{"user_id": data.UserID})
	})

	r := httptest.NewRequest("GET", "/authrequired", http.NoBody)

	tk, _ := security.CreateJWT("superuser")

	r.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tk,
	})

	rw := httptest.NewRecorder()

	router.ServeHTTP(rw, r)

	type response struct {
		UserID string `json:"user_id"`
	}

	res := new(response)

	err := json.Unmarshal(rw.Body.Bytes(), &res)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusAccepted, rw.Code)
	assert.Equal(t, "superuser", res.UserID)
}
