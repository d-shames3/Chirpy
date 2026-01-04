package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/d-shames3/chirpy/internal/database"
	"github.com/google/uuid"
)

type chirp struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type chirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

const (
	bleep    = "****"
	maxChars = 140
)

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	var chirp chirp
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&chirp); err != nil {
		respondWithError(w, http.StatusInternalServerError, "error reading chirp json")
		return
	}

	validChirp, err := validateChirp(chirp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	createChirpParams := database.CreateChirpParams{
		UserID: validChirp.UserID,
		Body:   validChirp.Body,
	}
	chirpData, err := cfg.db.CreateChirp(r.Context(), createChirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirpResponse := chirpResponse{
		ID:        chirpData.ID,
		UserID:    chirp.UserID,
		Body:      chirp.Body,
		CreatedAt: chirpData.CreatedAt,
		UpdatedAt: chirpData.UpdatedAt,
	}
	respondWithJSON(w, http.StatusCreated, chirpResponse)
}

func validateChirp(c chirp) (chirp, error) {
	if len(c.Body) > maxChars {
		return c, errors.New("Chirp is too long")
	}
	return stripProfanity(c), nil
}

func stripProfanity(c chirp) chirp {
	profanity := map[string]int{"kerfuffle": 0, "sharbert": 1, "fornax": 2}
	words := strings.Split(c.Body, " ")

	for i, word := range words {
		_, ok := profanity[strings.ToLower(word)]
		if ok {
			words[i] = bleep
		}
	}
	cleanWords := strings.Join(words, " ")

	cleanChirp := chirp{Body: cleanWords, UserID: c.UserID}
	return cleanChirp

}
