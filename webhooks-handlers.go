package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/phmshk/bootdev-http-servers/internal/auth"
)

func (cfg *apiConfig) upgradeUserRedHandler(w http.ResponseWriter, r *http.Request) {
	type data struct {
		UserID string `json:"user_id"`
	}

	type reqBody struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "StatusUnauthorized", err)

		return
	}

	if apiKey != cfg.polkaAPIKey {
		respondWithError(w, http.StatusUnauthorized, "StatusUnauthorized", err)
		return
	}

	var params reqBody
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "bad request body", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userUUID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing id", err)
		return
	}

	_, err = cfg.db.UpgradeUserToRed(r.Context(), userUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
