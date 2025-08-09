package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Bigesto/chirpy/handlers"
	"github.com/Bigesto/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("no url for the db, check if the .env exists")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(db)

	const port = "8080"
	const filepathRoot = "serverfiles"
	apiCfg := handlers.ApiConfig{
		Db: dbQueries,
	}

	mux := http.NewServeMux()
	httpHandler := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", httpHandler)))

	mux.HandleFunc("GET /api/healthz", handlers.HealthzHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetHandler)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.ValidateChirpsHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}
