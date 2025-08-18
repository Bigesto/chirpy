package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/Bigesto/chirpy/internal/auth"
	"github.com/Bigesto/chirpy/internal/database"
	"github.com/google/uuid"
)

var nonAllowedWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

type ChirpStruct struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *ApiConfig) CreateChirpsHandler(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		code := 401
		msg := "Invalid token."
		respondWithError(w, code, msg)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		code := 401
		msg := "Invalid token."
		respondWithError(w, code, msg)
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	bodyDecoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = bodyDecoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	validatedBody, err := validateChirp(params.Body)
	if err != nil {
		if err.Error() == "chirp is too long" {
			code := 400
			msg := "Chirp is too long"
			respondWithError(w, code, msg)
			return
		}
		code := 500
		msg := "Something went wrong during the validation."
		respondWithError(w, code, msg)
		return
	}

	chirpParams := database.CreateChirpParams{
		Body:   validatedBody,
		UserID: userID,
	}

	chirp, err := cfg.Db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		code := 500
		msg := "Something went wrong during the chirp's creation."
		respondWithError(w, code, msg)
		return
	}

	answer := ChirpStruct{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	code := 201
	respondWithJSON(w, code, answer)
}

func getCleanedBody(body string) string {
	slicedBody := strings.Split(body, " ")

	for i, word := range slicedBody {
		if _, bad := nonAllowedWords[strings.ToLower(word)]; bad {
			slicedBody[i] = "****"
		}
	}
	correctedBody := strings.Join(slicedBody, " ")

	return correctedBody
}

func validateChirp(body string) (string, error) {
	chirpLenght := len(body)
	if chirpLenght > 140 {
		err := errors.New("chirp is too long")
		return "", err
	}

	cleanedBody := getCleanedBody(body)

	return cleanedBody, nil
}

func (cfg *ApiConfig) GetAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	stringUserID := r.URL.Query().Get("author_id")
	sorted := r.URL.Query().Get("sort")

	var chirps []database.Chirp
	var err error

	if stringUserID == "" {
		chirps, err = cfg.Db.GetAllChirps(r.Context())
		if err != nil {
			code := 500
			msg := "Something went wrong getting the chirps."
			respondWithError(w, code, msg)
			return
		}
	} else {
		userID, err := uuid.Parse(stringUserID)
		if err != nil {
			code := 400
			msg := "Bad Request"
			respondWithError(w, code, msg)
			return
		}

		chirps, err = cfg.Db.GetChirpsByUserID(r.Context(), userID)
		if err != nil {
			code := 500
			msg := "Something went wrong getting the chirps."
			respondWithError(w, code, msg)
			return
		}
	}

	chirpsList := make([]ChirpStruct, 0, len(chirps))

	for _, chirp := range chirps {
		structuredChirp := ChirpStruct{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}

		chirpsList = append(chirpsList, structuredChirp)
	}

	if sorted == "desc" {
		sort.Slice(chirpsList, func(i, j int) bool {
			return chirpsList[i].CreatedAt.After(chirpsList[j].CreatedAt)
		})
	}

	code := 200
	respondWithJSON(w, code, chirpsList)
}

func (cfg *ApiConfig) GetChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		code := 500
		msg := "Something went wrong."
		respondWithError(w, code, msg)
		return
	}

	chirp, err := cfg.Db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			code := 404
			msg := "Chirp not found."
			respondWithError(w, code, msg)
			return
		}
		code := 500
		msg := "Something went wrong."
		respondWithError(w, code, msg)
		return
	}

	answer := ChirpStruct{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, 200, answer)
}

func (cfg *ApiConfig) DeleteChirpHandler(w http.ResponseWriter, r *http.Request) {
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

	chirpIDString := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		code := 400
		msg := "Bad request."
		respondWithError(w, code, msg)
		return
	}

	chirp, err := cfg.Db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			code := 404
			msg := "Chirp not found."
			respondWithError(w, code, msg)
			return
		}
		code := 500
		msg := "Something went wrong."
		respondWithError(w, code, msg)
		return
	}

	if chirp.UserID != userID {
		code := 403
		msg := "Non."
		respondWithError(w, code, msg)
		return
	}

	err = cfg.Db.DeleteChirpByID(r.Context(), chirpID)
	if err != nil {
		code := 500
		msg := "Something went wrong."
		respondWithError(w, code, msg)
		return
	}

	w.WriteHeader(204)
}
