package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/d-shames3/chirpy/internal/auth"
	"github.com/d-shames3/chirpy/internal/database"
	"github.com/google/uuid"
)

type userParams struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds *int   `json:"expires_in_seconds,omitempty"`
}

type userData struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token,omitempty"`
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	userParams := userParams{}
	defer r.Body.Close()
	decoder := json.NewDecoder((r.Body))
	if err := decoder.Decode(&userParams); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if userParams.ExpiresInSeconds == nil {
		defaultExpiresinSeconds := 3600
		userParams.ExpiresInSeconds = &defaultExpiresinSeconds
	} else if *userParams.ExpiresInSeconds > 3600 {
		defaultExpiresinSeconds := 3600
		userParams.ExpiresInSeconds = &defaultExpiresinSeconds
	}

	user, err := cfg.db.GetUser(r.Context(), userParams.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	match, err := auth.CheckPasswordHash(userParams.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	authToken, err := auth.MakeJWT(user.ID, cfg.serverSecret, time.Duration(*userParams.ExpiresInSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userData := userData{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     authToken,
	}
	respondWithJSON(w, http.StatusOK, userData)
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	userParams := userParams{}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userParams); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	hashedPassword, err := auth.HashPassword(userParams.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	createUserParams := database.CreateUserParams{
		Email:          userParams.Email,
		HashedPassword: hashedPassword,
	}

	user, err := cfg.db.CreateUser(r.Context(), createUserParams)
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
