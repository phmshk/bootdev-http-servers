package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/phmshk/bootdev-http-servers/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
	polkaAPIKey    string
}

func main() {
	const filepathRoot = "."
	const port = ":8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	mode := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaAPIKey := os.Getenv("POLKA_KEY")

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("An error occurred while accessing db: %s", err)
	}

	dbQueries := database.New(dbConn)
	apiConfig := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       mode,
		jwtSecret:      jwtSecret,
		polkaAPIKey:    polkaAPIKey,
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(filepathRoot))
	prefixHandler := http.StripPrefix("/app/", fs)

	mux.Handle("/app/", apiConfig.middlewareMetricsInc(prefixHandler))

	// API HANDLERS
	mux.HandleFunc("GET /api/healthz", firstHandler)

	mux.HandleFunc("GET /api/chirps", apiConfig.getChirpsHandler)
	mux.HandleFunc("POST /api/chirps", apiConfig.createChirpHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.getChirpByIDHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiConfig.deleteChirpHandler)

	mux.HandleFunc("POST /api/users", apiConfig.createUserHandler)
	mux.HandleFunc("PUT /api/users", apiConfig.updateUserCredentialsHandler)

	mux.HandleFunc("POST /api/login", apiConfig.userLoginHandler)
	mux.HandleFunc("POST /api/refresh", apiConfig.refreshAccessTokenHandler)
	mux.HandleFunc("POST /api/revoke", apiConfig.revokeTokenHandler)

	mux.HandleFunc("POST /api/polka/webhooks", apiConfig.upgradeUserRedHandler)

	// ADMIN HANDLERS
	mux.HandleFunc("GET /admin/metrics", apiConfig.requestsMetricsHandler)
	mux.HandleFunc("POST /admin/reset", apiConfig.resetMetricsHandler)

	server := &http.Server{
		Handler: mux,
		Addr:    port,
	}
	log.Println("Server starting on :8080...")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("An error occurred: %v", err)
	}
}
