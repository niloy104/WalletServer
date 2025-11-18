package walletB

import (
	"wallet/rest/middlewares"
	"wallet/wallet"
)

type Handler struct {
	svc         wallet.Service
	middlewares *middleware.Middlewares
}

func NewHandler(
	svc wallet.Service,
	middlewares *middleware.Middlewares,
) *Handler {
	return &Handler{
		svc:         svc,
		middlewares: middlewares,
	}
}
