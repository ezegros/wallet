package wallet

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ezegrosfeld/wallet/users/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type mockedRepo struct {
	mock.Mock
}

func (m *mockedRepo) Store(ctx context.Context, w *domain.Wallet) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func (m *mockedRepo) Find(ctx context.Context, user string) (*domain.Wallet, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*domain.Wallet), args.Error(1)
}

func TestCreate(t *testing.T) {
	repo := &mockedRepo{}

	repo.On("Find", mock.Anything, "user").Return(&domain.Wallet{}, ErrNotFound)
	repo.On("Store", mock.Anything, mock.Anything).Return(nil)

	l, _ := zap.NewDevelopment()

	sl := l.Sugar()

	s := NewService(sl, repo)

	ctx := context.Background()

	r := gin.Default()
	r.POST("/wallet", func(c *gin.Context) {
		c.JSON(200, gin.H{"address": "address", "seed": "seed123"})
	})

	go r.Run()

	os.Setenv("GENERATOR_URL", "http://127.0.0.1")

	w, err := s.Create(ctx, "user")
	assert.NoError(t, err)
	assert.NotNil(t, w)
}

func TestCreateWithConflictError(t *testing.T) {
	repo := &mockedRepo{}

	repo.On("Find", mock.Anything, "user").Return(&domain.Wallet{}, nil)

	l, _ := zap.NewDevelopment()

	sl := l.Sugar()

	s := NewService(sl, repo)

	ctx := context.Background()

	_, err := s.Create(ctx, "user")
	assert.Error(t, err)
	assert.Equal(t, ErrConflict, err)
}

func TestCreateWithErrorAtStore(t *testing.T) {
	repo := &mockedRepo{}

	repo.On("Find", mock.Anything, "user").Return(&domain.Wallet{}, ErrNotFound)
	repo.On("Store", mock.Anything, mock.Anything).Return(fmt.Errorf("error"))

	l, _ := zap.NewDevelopment()

	sl := l.Sugar()

	s := NewService(sl, repo)

	ctx := context.Background()

	r := gin.Default()
	r.POST("/wallet", func(c *gin.Context) {
		c.JSON(200, gin.H{"address": "address", "seed": "seed123"})
	})

	go r.Run()

	os.Setenv("GENERATOR_URL", "http://127.0.0.1")

	_, err := s.Create(ctx, "user")
	assert.Error(t, err)
}

func TestGet(t *testing.T) {
	repo := &mockedRepo{}

	repo.On("Find", mock.Anything, "user").Return(&domain.Wallet{}, nil)

	l, _ := zap.NewDevelopment()

	sl := l.Sugar()

	s := NewService(sl, repo)

	ctx := context.Background()

	w, err := s.Get(ctx, "user")
	assert.NoError(t, err)
	assert.NotNil(t, w)
}

func TestGetNotFound(t *testing.T) {
	repo := &mockedRepo{}

	repo.On("Find", mock.Anything, "user").Return(&domain.Wallet{}, ErrNotFound)

	l, _ := zap.NewDevelopment()

	sl := l.Sugar()

	s := NewService(sl, repo)

	ctx := context.Background()

	_, err := s.Get(ctx, "user")

	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
}

func TestGetAddress(t *testing.T) {
	repo := &mockedRepo{}

	repo.On("Find", mock.Anything, "user").Return(&domain.Wallet{}, nil)

	l, _ := zap.NewDevelopment()

	sl := l.Sugar()

	s := NewService(sl, repo)

	ctx := context.Background()

	r := gin.Default()
	r.POST("/wallet", func(c *gin.Context) {
		c.JSON(200, gin.H{"address": "address", "seed": "seed123"})
	})

	go r.Run()

	os.Setenv("GENERATOR_URL", "http://127.0.0.1")

	w, err := s.GetAddress(ctx, "user", 0)
	assert.NoError(t, err)
	assert.NotNil(t, w)
}

func TestGetAddressNotFound(t *testing.T) {
	repo := &mockedRepo{}

	repo.On("Find", mock.Anything, "user").Return(&domain.Wallet{}, ErrNotFound)

	l, _ := zap.NewDevelopment()

	sl := l.Sugar()

	s := NewService(sl, repo)

	ctx := context.Background()

	_, err := s.GetAddress(ctx, "user", 0)

	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
}
