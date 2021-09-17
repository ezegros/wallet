package wallet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/ezegrosfeld/wallet/users/internal/domain"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrConflict      = errors.New("user alredy exist")
	ErrNotFound      = errors.New("user not found")
	ErrWrongPassword = errors.New("the password is incorrect")
)

type Service interface {
	Create(ctx context.Context, userID string) (*domain.Wallet, error)
	Get(ctx context.Context, userID string) (*domain.Wallet, error)
	GetAddress(ctx context.Context, userID string, index int) (string, error)
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

// Create creates a new wallet for a given user
func (s *service) Create(ctx context.Context, userID string) (*domain.Wallet, error) {
	_, err := s.repo.Find(ctx, userID)
	if err == nil {
		return nil, ErrConflict
	}

	generatorURL := os.Getenv("GENERATOR_URL")
	resp, err := http.Post(fmt.Sprintf("%s/wallet", generatorURL), "application/json", http.NoBody)
	if err != nil {
		return nil, err
	}

	type response struct {
		Address string `json:"address"`
		Seed    string `json:"seed"`
	}

	var res response

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	wallet := &domain.Wallet{
		ID:     uuid.NewString(),
		UserID: userID,
		Seed:   res.Seed,
	}

	err = s.repo.Store(ctx, wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

// Get returns a wallet for a given user
func (s *service) Get(ctx context.Context, userID string) (*domain.Wallet, error) {
	wallet, err := s.repo.Find(ctx, userID)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

// GetAddress returns an address for a given user and index
func (s *service) GetAddress(ctx context.Context, userID string, index int) (string, error) {
	wallet, err := s.repo.Find(ctx, userID)
	if err != nil {
		return "", err
	}

	generatorURL := os.Getenv("GENERATOR_URL")

	resp, err := http.Get(fmt.Sprintf("%s/wallet?seed=%s&index=%d", generatorURL, wallet.Seed, index))
	if err != nil {
		return "", err
	}

	type response struct {
		Address string `json:"address"`
		Seed    string `json:"seed"`
	}

	var res response

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", err
	}

	return res.Address, nil
}
