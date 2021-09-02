package wallet

import (
	"fmt"
	"os"

	"github.com/ezegrosfeld/wallet/generator/internal/domain"
	"go.uber.org/zap"

	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	bip39 "github.com/tyler-smith/go-bip39"
)

type Service interface {
	Create() (domain.Wallet, error)
	Get(seedString string, index int) (domain.Wallet, error)
}

type service struct {
	log *zap.SugaredLogger
}

func NewService(log *zap.SugaredLogger) Service {
	return &service{
		log: log,
	}
}

func (s *service) Create() (domain.Wallet, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return domain.Wallet{}, err
	}

	mnemonic, _ := bip39.NewMnemonic(entropy)

	pwd := os.Getenv("WALLET_PASSWORD")

	seed := bip39.NewSeed(mnemonic, pwd)

	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		return domain.Wallet{}, err
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")

	account, err := wallet.Derive(path, false)
	if err != nil {
		return domain.Wallet{}, err
	}

	w := domain.Wallet{
		Seed:    mnemonic,
		Index:   0,
		Address: account.Address.Hex(),
	}

	return w, nil
}

func (s *service) Get(seedString string, index int) (domain.Wallet, error) {

	pwd := os.Getenv("WALLET_PASSWORD")

	seed := bip39.NewSeed(seedString, pwd)

	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		return domain.Wallet{}, err
	}

	path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", index))

	account, err := wallet.Derive(path, false)
	if err != nil {
		return domain.Wallet{}, err
	}

	w := domain.Wallet{
		Index:   index,
		Address: account.Address.Hex(),
	}

	return w, nil
}
