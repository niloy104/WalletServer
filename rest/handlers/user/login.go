package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"wallet/util"
)

type ReqLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req ReqLogin
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)

	if err != nil {
		fmt.Println(err)
		util.SendError(w, http.StatusBadRequest, "Invalid Request body")
		return
	}

	usr, err := h.svc.Find(req.Email, req.Password)

	if err != nil {
		util.SendError(w, http.StatusUnauthorized, "Unauthorized User")
		return
	}

	if usr == nil {
		util.SendError(w, http.StatusUnauthorized, "Unauthorized User")
		return
	}

	accessToken, err := util.CreateJwt(h.cnf.JwtSecretKey, util.Playload{
		Sub:       usr.ID,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Email:     usr.Email,
	})

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	util.SendData(w, http.StatusCreated, accessToken)

}
