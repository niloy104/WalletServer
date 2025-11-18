package repo

import (
	"context"
	"database/sql"
	"fmt"
	"wallet/domain"

	"github.com/jmoiron/sqlx"
)

type WalletRepository interface {
	GetByUserID(ctx context.Context, userID uint64) (*domain.Wallet, error)
	GetForUpdate(ctx context.Context, tx *sqlx.Tx, userID uint64) (*domain.Wallet, error)
	Create(ctx context.Context, wallet *domain.Wallet) error
	UpdateBalance(ctx context.Context, tx *sqlx.Tx, userID uint64, delta int64) error
}

type walletRepo struct {
	db *sqlx.DB
}

// Constructor
func NewWalletRepo(db *sqlx.DB) WalletRepository {
	return &walletRepo{db: db}
}

func (r *walletRepo) GetByUserID(ctx context.Context, userID uint64) (*domain.Wallet, error) {
	var w domain.Wallet
	query := `SELECT id, user_id, balance, currency, created_at, updated_at FROM wallets WHERE user_id=$1`
	err := r.db.GetContext(ctx, &w, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	return &w, nil
}

func (r *walletRepo) GetForUpdate(ctx context.Context, tx *sqlx.Tx, userID uint64) (*domain.Wallet, error) {
	var w domain.Wallet
	query := `SELECT id, user_id, balance, currency FROM wallets WHERE user_id=$1 FOR UPDATE`
	err := tx.GetContext(ctx, &w, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to lock wallet: %w", err)
	}
	return &w, nil
}

func (r *walletRepo) Create(ctx context.Context, wallet *domain.Wallet) error {
	query := `INSERT INTO wallets (user_id, balance, currency) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	return r.db.QueryRowxContext(ctx, query, wallet.UserID, wallet.Balance, wallet.Currency).
		Scan(&wallet.ID, &wallet.CreatedAt, &wallet.UpdatedAt)
}

func (r *walletRepo) UpdateBalance(ctx context.Context, tx *sqlx.Tx, userID uint64, delta int64) error {
	query := `UPDATE wallets SET balance = balance + $1, updated_at = NOW() WHERE user_id = $2`
	res, err := tx.ExecContext(ctx, query, delta, userID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("wallet not found for user %d", userID)
	}
	return nil
}
