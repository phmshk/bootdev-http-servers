package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"unicode/utf8"
)

func firstHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) requestsMetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	tmpl := `<html>
<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>
</html>`
	body := fmt.Sprintf(tmpl, cfg.fileserverHits.Load())
	w.Write([]byte(body))
}

func (cfg *apiConfig) resetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	body := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.Write([]byte(body))
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type params struct {
		Body string `json:"body"`
	}

	type successResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	const maxChirpLength = 140

	resp := params{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&resp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if utf8.RuneCountInString(resp.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedBody := getCleanedBody(resp.Body)

	respondWithJSON(w, http.StatusOK, successResponse{CleanedBody: cleanedBody})
}
