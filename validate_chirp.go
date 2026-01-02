package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type chirp struct {
	Body string `json:"body"`
}

const bleep = "****"

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirp := chirp{}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	validChirp := stripProfanity(chirp)

	respondWithJSON(w, http.StatusOK, map[string]string{"cleaned_body": validChirp.Body})
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
