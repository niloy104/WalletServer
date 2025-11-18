package walletB

import (
	"encoding/json"
	"net/http"
	"wallet/util"
)

type TransferRequest struct {
	ToUserID int64 `json:"to_user_id"`
	Amount   int64 `json:"amount"`
}

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	val := r.Context().Value("user_id")
	if val == nil {
		util.SendError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userIDInt, ok := val.(int)
	if !ok {
		util.SendError(w, http.StatusInternalServerError, "invalid user id")
		return
	}
	userID := uint64(userIDInt)

	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if req.Amount <= 0 {
		util.SendError(w, http.StatusBadRequest, "Amount must be greater than zero")
		return
	}

	if err := h.svc.Transfer(r.Context(), userID, uint64(req.ToUserID), req.Amount); err != nil {
		util.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.SendData(w, http.StatusOK, struct {
		Message  string `json:"message"`
		ToUserID int64  `json:"to_user_id"`
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
	}{
		Message:  "Transfer successful",
		ToUserID: req.ToUserID,
		Amount:   req.Amount,
		Currency: "USD",
	})
}
