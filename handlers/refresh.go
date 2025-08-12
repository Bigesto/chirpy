package handlers

import (
	"net/http"
	"time"

	"github.com/Bigesto/chirpy/internal/auth"
)

func (cfg *ApiConfig) RefreshHandler(w http.ResponseWriter, r *http.Request) {
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

	userID, err := cfg.Db.GetUserIDByRefreshToken(r.Context(), checkedToken.Token)
	if err != nil {
		code := 401
		msg := "token invalid"
		respondWithError(w, code, msg)
		return
	}

	accessToken, err := auth.MakeJWT(userID, cfg.Secret, time.Duration(3600)*time.Second)
	if err != nil {
		code := 500
		msg := "token invalid"
		respondWithError(w, code, msg)
		return
	}

	type TokenResponse struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, 200, TokenResponse{Token: accessToken})
}
