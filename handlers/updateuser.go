package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Bigesto/chirpy/internal/auth"
	"github.com/Bigesto/chirpy/internal/database"
)

func (cfg *ApiConfig) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		code := 401
		msg := "token invalid"
		respondWithError(w, code, msg)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.Secret)
	if err != nil {
		code := 401
		msg := "Invalid token."
		respondWithError(w, code, msg)
		return
	}

	type parameters struct {
		NewPassword string `json:"password"`
		NewEmail    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	hashedPassword, err := auth.HashPassword(params.NewPassword)
	if err != nil {
		code := 500
		msg := "Internal error"
		respondWithError(w, code, msg)
		return
	}

	userParams := database.UpdateUserParams{
		ID:             userID,
		Email:          params.NewEmail,
		HashedPassword: hashedPassword,
	}

	err = cfg.Db.UpdateUser(r.Context(), userParams)
	if err != nil {
		code := 500
		msg := "Update error"
		respondWithError(w, code, msg)
		return
	}

	user, err := cfg.Db.GetUserByID(r.Context(), userID)
	if err != nil {
		code := 500
		msg := "User error"
		respondWithError(w, code, msg)
		return
	}

	answer := User{
		ID:          userID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJSON(w, 200, answer)
}
