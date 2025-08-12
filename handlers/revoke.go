package handlers

import (
	"net/http"

	"github.com/Bigesto/chirpy/internal/auth"
)

func (cfg *ApiConfig) RevokeHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		code := 400
		msg := "token invalid"
		respondWithError(w, code, msg)
		return
	}

	checkedToken, err := cfg.Db.CheckRefreshToken(r.Context(), refreshToken)
	if err != nil {
		code := 401
		msg := "token invalid"
		respondWithError(w, code, msg)
		return
	}

	err = cfg.Db.RevokeRefreshToken(r.Context(), checkedToken.Token)
	if err != nil {
		code := 401
		msg := "token invalid"
		respondWithError(w, code, msg)
		return
	}

	w.WriteHeader(204)
}
