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
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	polkaKey := os.Getenv("POLKA_KEY")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(db)

	const port = "8080"
	const filepathRoot = "serverfiles"
	apiCfg := handlers.ApiConfig{
		Db:       dbQueries,
		Platform: platform,
		Secret:   secret,
		PolkaKey: polkaKey,
	}

	mux := http.NewServeMux()
	httpHandler := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", httpHandler)))

	mux.HandleFunc("GET /api/healthz", handlers.HealthzHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetHandler)
	mux.HandleFunc("POST /api/users", apiCfg.CreateUserHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.CreateChirpsHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.GetAllChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.GetChirpByIDHandler)
	mux.HandleFunc("POST /api/login", apiCfg.LoginHandler)
	mux.HandleFunc("POST /api/refresh", apiCfg.RefreshHandler)
	mux.HandleFunc("POST /api/revoke", apiCfg.RevokeHandler)
	mux.HandleFunc("PUT /api/users", apiCfg.UpdateUserHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.DeleteChirpHandler)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.UpgradeUserHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}
