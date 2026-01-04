package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type userEmail struct {
	Email string `json:"email"`
}

type userData struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	userEmail := userEmail{}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userEmail); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), userEmail.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	userData := userData{}
	if err := json.Unmarshal(jsonUser, &userData); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusCreated, userData)
}
