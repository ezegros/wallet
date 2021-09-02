package wallet

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// Test the correct generation of a wallet
func TestWalletCreation(t *testing.T) {
	service := NewService(&zap.SugaredLogger{})

	wallet, err := service.Create()
	assert.NoError(t, err)
	assert.NotNil(t, wallet.Address)
	assert.NotNil(t, wallet.Seed)

	// Check that the seed has exactly 12 words
	assert.Equal(t, 12, len(strings.Split(wallet.Seed, " ")))

	// Check uniqueness of the address and seed
	wallet2, err := service.Create()
	assert.NoError(t, err)

	assert.NotEqual(t, wallet.Address, wallet2.Address)
	assert.NotEqual(t, wallet.Seed, wallet2.Seed)

}

// Test the corret return of a wallet from a seed
func TestWalletGet(t *testing.T) {
	service := NewService(&zap.SugaredLogger{})

	wallet, err := service.Create()
	assert.NoError(t, err)

	wallet2, err := service.Get(wallet.Seed, 0)
	assert.NoError(t, err)

	assert.Equal(t, wallet.Address, wallet2.Address)
}
