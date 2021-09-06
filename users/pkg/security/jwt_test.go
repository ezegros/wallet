package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWTCreation(t *testing.T) {
	token, err := CreateJWT("mysuperuser")
	assert.NoError(t, err)

	assert.NotNil(t, token)
}

func TestJWT(t *testing.T) {
	token, _ := CreateJWT("mysuperuser")

	validatedToken, err := ValidateJWT(token)

	assert.NoError(t, err)
	assert.Equal(t, "mysuperuser", validatedToken.UserID)
}
