package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

var nonAllowedWords = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

func (cfg *ApiConfig) ValidateChirpsHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	chirpLenght := len(params.Body)
	log.Printf("Chirp lenght is: %v", chirpLenght)
	if chirpLenght > 140 {
		msg := "Chirp is too long"
		code := 400
		respondWithError(w, code, msg)
		return
	}

	type respBodyCleaned struct {
		CleanedBody string `json:"cleaned_body"`
	}

	slicedBody := strings.Split(params.Body, " ")

	for i, word := range slicedBody {
		for _, swearWord := range nonAllowedWords {
			if strings.ToLower(word) == swearWord {
				slicedBody[i] = "****"
				break
			}
		}
	}

	correctedBody := strings.Join(slicedBody, " ")
	code := 200
	response := respBodyCleaned{CleanedBody: correctedBody}
	respondWithJSON(w, code, response)
}
