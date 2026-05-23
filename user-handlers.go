package main

import (
	"encoding/json"
	"net/http"

	"github.com/phmshk/bootdev-http-servers/internal/auth"
	"github.com/phmshk/bootdev-http-servers/internal/database"
)

type UserLoginStruct struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody UserLoginStruct

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding json", err)
		return
	}

	hashedPassword, err := auth.HashPassword(reqBody.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          reqBody.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (cfg *apiConfig) updateUserCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Wrong token format provided", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token is not valid", err)
		return
	}

	var params requestBody
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding request body", err)
		return
	}

	newHashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing password", err)
		return
	}

	updatedUser, err := cfg.db.UpdateUserCredentials(r.Context(), database.UpdateUserCredentialsParams{
		Email:          params.Email,
		HashedPassword: newHashedPassword,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating user data", err)
		return
	}

	respondWithJSON(w, http.StatusOK, UserResponse{
		Email:       updatedUser.Email,
		ID:          updatedUser.ID,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
		IsChirpyRed: updatedUser.IsChirpyRed,
	})
}
