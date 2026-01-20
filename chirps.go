package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/d-shames3/chirpy/internal/auth"
	"github.com/d-shames3/chirpy/internal/database"
	"github.com/google/uuid"
)

type chirp struct {
	Body string `json:"body"`
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

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpIdPath := r.PathValue("chirpId")
	if chirpIdPath == "" {
		respondWithError(w, http.StatusBadRequest, "no chirp id provided")
		return
	}

	chirpId, err := uuid.Parse(chirpIdPath)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "chirp id is not in UUID format")
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	chirpResponse := chirpResponse{
		ID:        chirp.ID,
		UserID:    chirp.UserID,
		Body:      chirp.Body,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	}

	respondWithJSON(w, http.StatusOK, chirpResponse)
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var chirpResponses []chirpResponse
	for _, chirp := range chirps {
		chirpResponse := chirpResponse{
			ID:        chirp.ID,
			UserID:    chirp.UserID,
			Body:      chirp.Body,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
		}
		chirpResponses = append(chirpResponses, chirpResponse)
	}

	respondWithJSON(w, http.StatusOK, chirpResponses)
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	var chirp chirp
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&chirp); err != nil {
		respondWithError(w, http.StatusInternalServerError, "error reading chirp json")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.serverSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	validChirp, err := validateChirp(chirp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	createChirpParams := database.CreateChirpParams{
		UserID: userId,
		Body:   validChirp.Body,
	}
	chirpData, err := cfg.db.CreateChirp(r.Context(), createChirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirpResponse := chirpResponse{
		ID:        chirpData.ID,
		UserID:    userId,
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

	cleanChirp := chirp{Body: cleanWords}
	return cleanChirp

}
