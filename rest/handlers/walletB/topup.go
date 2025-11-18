package walletB

import (
	"encoding/json"
	"net/http"
	"wallet/util"
)

type TopUpRequest struct {
	Amount int64 `json:"amount"`
}

func (h *Handler) TopUp(w http.ResponseWriter, r *http.Request) {
	val := r.Context().Value("user_id")
	if val == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var userID uint64
	switch v := val.(type) {
	case int:
		userID = uint64(v)
	case int64:
		userID = uint64(v)
	case float64:
		userID = uint64(v)
	default:
		http.Error(w, "invalid user id", http.StatusInternalServerError)
		return
	}

	var req TopUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.Amount <= 0 {
		util.SendError(w, http.StatusBadRequest, "Amount must be greater than zero")
		return
	}

	if err := h.svc.TopUp(r.Context(), userID, req.Amount); err != nil {
		util.SendError(w, http.StatusInternalServerError, "Failed to top up")
		return
	}

	util.SendData(w, http.StatusOK, struct {
		Message  string `json:"message"`
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
	}{
		Message:  "Top up successful",
		Amount:   req.Amount,
		Currency: "USD",
	})
}
