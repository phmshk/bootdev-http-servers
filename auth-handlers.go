package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/phmshk/bootdev-http-servers/internal/auth"
	"github.com/phmshk/bootdev-http-servers/internal/database"
)

func (cfg *apiConfig) userLoginHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody UserLoginStruct
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding json", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), reqBody.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(reqBody.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error checking password", err)
		return
	}

	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT token", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error saving refresh token to database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, UserResponseWithToken{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        accessToken,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	})
}

func (cfg *apiConfig) refreshAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshTokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid token format", err)
		return
	}

	dbToken, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshTokenStr)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	if dbToken.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired", nil)
		return
	}

	if dbToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token revoked", nil)
		return
	}

	accessToken, err := auth.MakeJWT(dbToken.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token", err)
		return
	}

	type refreshResponse struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, refreshResponse{
		Token: accessToken,
	})
}

func (cfg *apiConfig) revokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshTokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid token format", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshTokenStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
