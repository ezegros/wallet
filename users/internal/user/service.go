package user

import (
	"context"
	"errors"

	"github.com/ezegrosfeld/wallet/users/internal/domain"
	"github.com/ezegrosfeld/wallet/users/pkg/security"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrConflict      = errors.New("user alredy exist")
	ErrNotFound      = errors.New("user not found")
	ErrWrongPassword = errors.New("the password is incorrect")
)

type Service interface {
	Create(ctx context.Context, username, password string) (*domain.User, error)
	Login(ctx context.Context, username, password string) (*domain.User, error)
}

type service struct {
	repo Repository
	log  *zap.SugaredLogger
}

// NewService returns a service struct which must implement the Service interface
func NewService(log *zap.SugaredLogger, repo Repository) Service {
	return &service{
		log:  log.Named("User Service"),
		repo: repo,
	}
}

// Create creates a new user using the repo and receiveng the data from the handler
func (s *service) Create(ctx context.Context, username, password string) (*domain.User, error) {
	// Check if user exists
	exists := s.repo.Exists(username)
	if exists {
		return nil, ErrConflict
	}

	// Hash the password
	pwd, _ := security.HashPassword(password)

	// Map the values
	user := &domain.User{
		ID:       uuid.New().String(),
		Username: username,
		Password: pwd,
	}

	// Store the new user in the table
	err := s.repo.Store(ctx, user)
	if err != nil {
		return nil, err
	}

	// Prevent the password from being returned back to the user
	user.Password = ""

	// Log the creation
	s.log.Infow("Created user", "username", username)

	return user, nil
}

func (s *service) Login(ctx context.Context, username, password string) (*domain.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, ErrNotFound
	}

	err = security.CompareHashAndPassword(user.Password, password)
	if err != nil {
		return nil, ErrWrongPassword
	}

	return user, nil
}
