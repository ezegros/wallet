package user

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ezegrosfeld/wallet/users/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type mockedRepo struct {
	mock.Mock
}

func (r *mockedRepo) Store(ctx context.Context, user *domain.User) error {
	args := r.Called(user)
	return args.Error(0)
}

func (r *mockedRepo) Exists(username string) bool {
	args := r.Called(username)

	return args.Bool(0)
}

func TestServiceCreate(t *testing.T) {
	mr := &mockedRepo{}

	l, _ := zap.NewProduction()

	user := &domain.User{
		Username: "username",
		Password: "superstrongpassword",
	}

	mr.On("Store", mock.AnythingOfType("*domain.User")).Return(nil)
	mr.On("Exists", user.Username).Return(false)

	s := NewService(l.Sugar(), mr)

	su, err := s.Create(context.Background(), user.Username, user.Username)

	assert.NoError(t, err)
	assert.Equal(t, user.Username, su.Username)
	assert.NotEmpty(t, su.ID)
	assert.Empty(t, su.Password)
}

func TestServiceCreateError(t *testing.T) {
	mr := &mockedRepo{}

	l, _ := zap.NewProduction()

	user := &domain.User{
		Username: "username",
		Password: "superstrongpassword",
	}

	mr.On("Store", mock.AnythingOfType("*domain.User")).Return(fmt.Errorf("Something wrong happened"))
	mr.On("Exists", user.Username).Return(false)

	s := NewService(l.Sugar(), mr)

	_, err := s.Create(context.Background(), user.Username, user.Username)

	assert.Error(t, err)
}

func TestServiceCreateErrorAlredyExists(t *testing.T) {
	mr := &mockedRepo{}

	l, _ := zap.NewProduction()

	user := &domain.User{
		Username: "username",
		Password: "superstrongpassword",
	}

	mr.On("Store", mock.AnythingOfType("*domain.User")).Return(nil)
	mr.On("Exists", user.Username).Return(true)

	s := NewService(l.Sugar(), mr)

	_, err := s.Create(context.Background(), user.Username, user.Username)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrConflict))
}
