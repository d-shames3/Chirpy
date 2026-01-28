package main

import (
	"encoding/json"
	"net/http"

	"github.com/d-shames3/chirpy/internal/auth"
	"github.com/google/uuid"
)

type polkaPaidUserData struct {
	Event   string               `json:"event"`
	Payload polkaPaidUserPayload `json:"data"`
}

type polkaPaidUserPayload struct {
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) polkaPaidUserWebhookHandler(w http.ResponseWriter, r *http.Request) {
	polkaKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if cfg.apiKey != polkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var rawWebhook polkaPaidUserData
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&rawWebhook); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if rawWebhook.Event != "user.upgraded" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.db.UpdateUserChirpyRedStatus(r.Context(), rawWebhook.Payload.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
