package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Bigesto/chirpy/internal/auth"
	"github.com/Bigesto/chirpy/internal/database"
)

func (cfg *ApiConfig) LoginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	user, err := cfg.Db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			code := 401
			msg := "Incorrect email or password"
			respondWithError(w, code, msg)
			return
		}
		code := 500
		msg := "Something went wrong."
		respondWithError(w, code, msg)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		code := 401
		msg := "Incorrect email or password"
		respondWithError(w, code, msg)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.Secret, time.Duration(3600)*time.Second)
	if err != nil {
		code := 500
		msg := "token invalid"
		respondWithError(w, code, msg)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		code := 500
		msg := "token invalid"
		respondWithError(w, code, msg)
		return
	}

	refreshParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	}

	err = cfg.Db.CreateRefreshToken(r.Context(), refreshParams)
	if err != nil {
		code := 500
		msg := "token invalid"
		respondWithError(w, code, msg)
		return
	}

	answer := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed,
		Token:        token,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, 200, answer)
}
