package wallet

import (
	"context"
	"errors"
	"wallet/domain"
	"wallet/repo"

	"github.com/jmoiron/sqlx"
)

// Service interface
type Service interface {
	GetBalance(ctx context.Context, userID uint64) (*domain.Wallet, error)
	TopUp(ctx context.Context, userID uint64, amount int64) error
	Transfer(ctx context.Context, fromUserID, toUserID uint64, amount int64) error
}

// Concrete service
type walletService struct {
	walletRepo repo.WalletRepository
	db         repo.DBExecutor
}

// Constructor
func NewService(walletRepo repo.WalletRepository, db repo.DBExecutor) Service {
	return &walletService{
		walletRepo: walletRepo,
		db:         db,
	}
}

// GetBalance returns user's wallet
func (s *walletService) GetBalance(ctx context.Context, userID uint64) (*domain.Wallet, error) {
	return s.walletRepo.GetByUserID(ctx, userID)
}

// TopUp adds funds to a wallet
func (s *walletService) TopUp(ctx context.Context, userID uint64, amount int64) error {
	if amount <= 0 {
		return errors.New("invalid amount") // fallback if domain.ErrInvalidAmount missing
	}

	return s.db.WithTx(ctx, func(tx *sqlx.Tx) error {
		wallet, err := s.walletRepo.GetForUpdate(ctx, tx, userID)
		if err != nil {
			return err
		}

		if wallet == nil {
			wallet = &domain.Wallet{UserID: userID, Balance: 0, Currency: "USD"}
			if err := s.walletRepo.Create(ctx, wallet); err != nil {
				return err
			}
		}

		return s.walletRepo.UpdateBalance(ctx, tx, userID, amount)
	})
}

// Transfer moves funds between wallets
func (s *walletService) Transfer(ctx context.Context, fromUserID, toUserID uint64, amount int64) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}
	if fromUserID == toUserID {
		return errors.New("cannot transfer to self")
	}

	return s.db.WithTx(ctx, func(tx *sqlx.Tx) error {
		// Lock sender
		sender, err := s.walletRepo.GetForUpdate(ctx, tx, fromUserID)
		if err != nil {
			return err
		}
		if sender == nil || sender.Balance < amount {
			return errors.New("insufficient balance")
		}

		// Lock receiver
		receiver, _ := s.walletRepo.GetForUpdate(ctx, tx, toUserID)
		if receiver == nil {
			receiver = &domain.Wallet{UserID: toUserID, Balance: 0, Currency: "USD"}
			if err := s.walletRepo.Create(ctx, receiver); err != nil {
				return err
			}
		}

		if err := s.walletRepo.UpdateBalance(ctx, tx, fromUserID, -amount); err != nil {
			return err
		}
		return s.walletRepo.UpdateBalance(ctx, tx, toUserID, amount)
	})
}
