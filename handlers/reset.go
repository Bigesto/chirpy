package handlers

import (
	"log"
	"net/http"
)

func (cfg *ApiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	log.Println("Hits counter reset.")
	w.Write([]byte("Hits reset to 0"))
}
