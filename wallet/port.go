package wallet

import (
	"wallet/domain"

	"github.com/jmoiron/sqlx"
)

type WalletRepo interface {
	GetByUserID(userID uint64) (*domain.Wallet, error)
	GetForUpdate(tx *sqlx.Tx, userID uint64) (*domain.Wallet, error)
	Create(wallet *domain.Wallet) error
	UpdateBalance(tx *sqlx.Tx, userID uint64, delta int64) error
}
