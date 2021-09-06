package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHashing(t *testing.T) {
	s, err := HashPassword("passord")
	assert.NoError(t, err)
	assert.NotEqual(t, "password", s)
}

func TestCompareHashPassword(t *testing.T) {
	s, _ := HashPassword("password")

	err := CompareHashAndPassword(s, "password")
	assert.NoError(t, err)
}

func TestCompareHashPasswordFail(t *testing.T) {
	s, _ := HashPassword("password")

	err := CompareHashAndPassword(s, "nosamepassword")
	assert.Error(t, err)
}
