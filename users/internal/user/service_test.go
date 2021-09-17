package user

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ezegrosfeld/wallet/users/internal/domain"
	"github.com/ezegrosfeld/wallet/users/pkg/security"
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

func (r *mockedRepo) GetByUsername(username string) (*domain.User, error) {
	args := r.Called(username)
	return args.Get(0).(*domain.User), args.Error(1)
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

func TestLogin(t *testing.T) {
	mr := &mockedRepo{}

	l, _ := zap.NewProduction()

	user := &domain.User{
		Username: "username",
		Password: "superstrongpassword",
	}

	hp, _ := security.HashPassword(user.Password)

	suser := &domain.User{
		Username: "username",
		Password: hp,
	}

	mr.On("GetByUsername", user.Username).Return(suser, nil)

	s := NewService(l.Sugar(), mr)

	usr, err := s.Login(context.Background(), user.Username, user.Password)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, usr.Username)
	assert.Equal(t, user.ID, usr.ID)
}

func TestFailedLoginNotFound(t *testing.T) {
	mr := &mockedRepo{}

	l, _ := zap.NewProduction()

	user := &domain.User{
		Username: "username",
		Password: "superstrongpassword",
	}

	mr.On("GetByUsername", user.Username).Return(user, fmt.Errorf("not found"))

	s := NewService(l.Sugar(), mr)

	_, err := s.Login(context.Background(), user.Username, user.Password)
	assert.Error(t, err)
	assert.Equal(t, err, ErrNotFound)
}

func TestLoginPasswordIncorrect(t *testing.T) {
	mr := &mockedRepo{}

	l, _ := zap.NewProduction()

	user := &domain.User{
		Username: "username",
		Password: "superstrongpassword",
	}

	hp, _ := security.HashPassword("anotherpassword")

	suser := &domain.User{
		Username: "username",
		Password: hp,
	}

	mr.On("GetByUsername", user.Username).Return(suser, nil)

	s := NewService(l.Sugar(), mr)

	_, err := s.Login(context.Background(), user.Username, user.Password)
	assert.Error(t, err)
	assert.Equal(t, err, ErrWrongPassword)
}
