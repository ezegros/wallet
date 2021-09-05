package user

import (
	"context"

	"github.com/ezegrosfeld/wallet/users/internal/domain"
	"github.com/ezegrosfeld/wallet/users/pkg/security"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, username, password string) (*domain.User, error)
}

type service struct {
	repo Repository
	log  *zap.SugaredLogger
}

func NewService(log *zap.SugaredLogger, repo Repository) Service {
	return &service{
		log:  log.Named("User Service"),
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, username, password string) (*domain.User, error) {
	pwd, _ := security.HashPassword(password)

	user := &domain.User{
		ID:       uuid.New().String(),
		Username: username,
		Password: pwd,
	}

	err := s.repo.Store(ctx, user)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	s.log.Infow("Created user", "username", username)

	return user, nil
}
