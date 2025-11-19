package wallet

import (
	"context"
	"errors"
	"wallet/domain"
	"wallet/repo"

	"github.com/jmoiron/sqlx"
)

type Service interface {
	GetBalance(ctx context.Context, userID uint64) (*domain.Wallet, error)
	TopUp(ctx context.Context, userID uint64, amount int64) error
	Transfer(ctx context.Context, fromUserID, toUserID uint64, amount int64) error
	GetTransactionHistory(ctx context.Context, userID uint64) ([]domain.Transaction, error)
}

type walletService struct {
	walletRepo repo.WalletRepository
	db         repo.DBExecutor
}

func NewService(walletRepo repo.WalletRepository, db repo.DBExecutor) Service {
	return &walletService{
		walletRepo: walletRepo,
		db:         db,
	}
}

func (s *walletService) GetBalance(ctx context.Context, userID uint64) (*domain.Wallet, error) {
	return s.walletRepo.GetByUserID(ctx, userID)
}

func (s *walletService) TopUp(ctx context.Context, userID uint64, amount int64) error {
	if amount <= 0 {
		return errors.New("invalid amount")
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

		if err := s.walletRepo.UpdateBalance(ctx, tx, userID, amount); err != nil {
			return err
		}

		txn := &domain.Transaction{
			FromUserID: nil,
			ToUserID:   &userID,
			Amount:     amount,
			Type:       "topup",
			Status:     "success",
		}

		return s.walletRepo.CreateTransaction(ctx, tx, txn)
	})
}

func (s *walletService) Transfer(ctx context.Context, fromUserID, toUserID uint64, amount int64) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}
	if fromUserID == toUserID {
		return errors.New("cannot transfer to self")
	}

	return s.db.WithTx(ctx, func(tx *sqlx.Tx) error {

		sender, err := s.walletRepo.GetForUpdate(ctx, tx, fromUserID)
		if err != nil {
			return err
		}
		if sender == nil || sender.Balance < amount {
			return errors.New("insufficient balance")
		}

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
		if err := s.walletRepo.UpdateBalance(ctx, tx, toUserID, amount); err != nil {
			return err
		}

		txnOut := &domain.Transaction{
			FromUserID: &fromUserID,
			ToUserID:   &toUserID,
			Amount:     amount,
			Type:       "transfer_out",
			Status:     "success",
		}
		if err := s.walletRepo.CreateTransaction(ctx, tx, txnOut); err != nil {
			return err
		}

		txnIn := &domain.Transaction{
			FromUserID: &fromUserID,
			ToUserID:   &toUserID,
			Amount:     amount,
			Type:       "transfer_in",
			Status:     "success",
		}
		return s.walletRepo.CreateTransaction(ctx, tx, txnIn)
	})
}

func (s *walletService) GetTransactionHistory(ctx context.Context, userID uint64) ([]domain.Transaction, error) {
	return s.walletRepo.GetTransactionsByUserID(ctx, userID)
}
