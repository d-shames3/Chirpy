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
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userData struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

const defaultExpiresinSeconds = 3600

func (cfg *apiConfig) revokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	_, err = cfg.db.RevokeToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	refreshTokenData, err := cfg.db.GetToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if refreshTokenData.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token has been revoked")
		return
	}

	if time.Now().UTC().After(refreshTokenData.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token is expired")
		return
	}

	authToken, err := auth.MakeJWT(refreshTokenData.UserID, cfg.serverSecret, time.Duration(defaultExpiresinSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := struct {
		Token string `json:"token"`
	}{Token: authToken}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	userParams := userParams{}
	defer r.Body.Close()
	decoder := json.NewDecoder((r.Body))
	if err := decoder.Decode(&userParams); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
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

	authToken, err := auth.MakeJWT(user.ID, cfg.serverSecret, time.Duration(defaultExpiresinSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	createTokenParams := database.CreateTokenParams{
		UserID: user.ID,
		Token:  refreshToken,
	}

	refreshTokenData, err := cfg.db.CreateToken(r.Context(), createTokenParams)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userData := userData{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        authToken,
		RefreshToken: refreshTokenData.Token,
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
