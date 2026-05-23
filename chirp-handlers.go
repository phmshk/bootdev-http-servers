package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/phmshk/bootdev-http-servers/internal/auth"
	"github.com/phmshk/bootdev-http-servers/internal/database"
)

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Body string `json:"body"`
	}

	const maxChirpLength = 140

	resp := params{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&resp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	gotToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid header format", err)
		return
	}

	userID, err := auth.ValidateJWT(gotToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	if utf8.RuneCountInString(resp.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedBody := getCleanedBody(resp.Body)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	var chirps []database.Chirp
	var err error
	authorIDStr := r.URL.Query().Get("author_id")
	sortParam := r.URL.Query().Get("sort")

	if authorIDStr != "" {
		authorUUID, err := uuid.Parse(authorIDStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID format", err)
			return
		}
		chirps, err = cfg.db.GetChirpsByAuthor(r.Context(), authorUUID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps for author", err)
			return
		}
	} else {
		chirps, err = cfg.db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
			return
		}
	}

	type chirpResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	response := []chirpResponse{}
	for _, chirp := range chirps {
		response = append(response, chirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	if sortParam == "desc" {
		slices.SortFunc(response, func(a, b chirpResponse) int {
			if a.CreatedAt.After(b.CreatedAt) {
				return -1
			}
			if a.CreatedAt.Before(b.CreatedAt) {
				return 1
			}
			return 0
		})
	} else {
		slices.SortFunc(response, func(a, b chirpResponse) int {
			if a.CreatedAt.Before(b.CreatedAt) {
				return -1
			}
			if a.CreatedAt.After(b.CreatedAt) {
				return 1
			}
			return 0
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) getChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	chirpIDVal := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDVal)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid chirp Id", err)
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}
	userID, err := auth.ValidateJWT(tokenStr, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	chirpUUID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error converting id", err)
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
		return
	}
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "operation forbidden", nil)
		return
	}

	err = cfg.db.DeleteChirpByID(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error deleting chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
