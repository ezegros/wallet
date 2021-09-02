package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainProgram(t *testing.T) {
	go main()

	r, err := http.Get("http://localhost:8080/health/")
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, r.StatusCode)
}
