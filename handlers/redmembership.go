package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Bigesto/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) UpgradeUserHandler(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		code := 401
		msg := "Access denied"
		respondWithError(w, code, msg)
	}

	if apiKey != cfg.PolkaKey {
		code := 401
		msg := "Access denied"
		respondWithError(w, code, msg)
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	stringUID := params.Data.UserId
	userID, err := uuid.Parse(stringUID)
	if err != nil {
		code := 400
		msg := "Bad request."
		respondWithError(w, code, msg)
		return
	}

	err = cfg.Db.UpgradeToRed(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			code := 404
			msg := "User not found."
			respondWithError(w, code, msg)
			return
		}
		code := 500
		msg := "Something went wrong."
		respondWithError(w, code, msg)
		return
	}

	w.WriteHeader(204)
}
