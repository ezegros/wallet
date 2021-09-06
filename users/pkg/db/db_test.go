package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseInitialize(t *testing.T) {
	InitializeDatabase()

	db := DynamoDB

	assert.NotNil(t, db)
}
