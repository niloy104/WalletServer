package walletB

import (
	"context"
	"wallet/domain"
)

type Service interface {
	GetBalance(ctx context.Context, userID uint64) (*domain.Wallet, error)
	TopUp(ctx context.Context, userID uint64, amount int64) error
	Transfer(ctx context.Context, fromUserID, toUserID uint64, amount int64) error
}

type ReqTopUp struct {
	Amount int64 `json:"amount"`
}

type ReqTransfer struct {
	ToUserID uint64 `json:"to_user_id"`
	Amount   int64  `json:"amount"`
}

type ResBalance struct {
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}
