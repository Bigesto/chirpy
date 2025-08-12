package handlers

import (
	"log"
	"net/http"
)

func (cfg *ApiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		code := 403
		msg := "403 Forbidden"
		respondWithError(w, code, msg)
		return
	}

	err := cfg.Db.DeleteAllUsers(r.Context())
	if err != nil {
		log.Printf("Error deleting users: %s", err)
		code := 500
		msg := "something went wrong during the deletion"
		respondWithError(w, code, msg)
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	log.Println("Database reset. Hits counter reset.")
	w.Write([]byte("Users database cleaned. Hits reset to 0"))
}
