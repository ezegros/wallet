package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMainProgram(t *testing.T) {
	router := gin.Default()

	healthCheck(router)

	r := httptest.NewRequest("GET", "/health", nil)

	rw := httptest.NewRecorder()

	router.ServeHTTP(rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)
}
