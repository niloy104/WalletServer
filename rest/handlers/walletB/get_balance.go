package walletB

import (
	"net/http"
	"wallet/util"
)

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {

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

	wallet, err := h.svc.GetBalance(r.Context(), userID)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if wallet == nil {
		util.SendData(w, http.StatusOK, struct {
			Balance  int64  `json:"balance"`
			Currency string `json:"currency"`
		}{
			Balance:  0,
			Currency: "USD",
		})
		return
	}

	util.SendData(w, http.StatusOK, struct {
		UserID   uint64 `json:"user_id"`
		Balance  int64  `json:"balance"`
		Currency string `json:"currency"`
	}{
		UserID:   wallet.UserID,
		Balance:  wallet.Balance,
		Currency: wallet.Currency,
	})

}
