package domain

import "errors"

type Wallet struct {
	ID        uint64 `db:"id"`
	UserID    uint64 `db:"user_id"`
	Balance   int64  `db:"balance"`
	Currency  string `db:"currency"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

// domain errors
var (
	ErrInvalidAmount       = errors.New("amount must be greater than zero")
	ErrTransferToSelf      = errors.New("cannot transfer to self")
	ErrInsufficientBalance = errors.New("insufficient balance")
)
