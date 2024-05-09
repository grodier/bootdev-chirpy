package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

func (cfg *apiConfig) handlerPolka(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header is required")
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "ApiKey" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header must be in the format 'ApikKey {key}'")
		return
	}

	apiKey := parts[1]
	if apiKey != os.Getenv("POLKA_API_KEY") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	type Data struct {
		UserId int `json:"user_id"`
	}

	type parameters struct {
		Event    string `json:"event"`
		DataItem Data   `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil || params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusOK)
		return
	}

	user, err := cfg.DB.GetUser(params.DataItem.UserId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	user.IsChirpyRed = true

	cfg.DB.UpdateUser(params.DataItem.UserId, user)

	w.WriteHeader(http.StatusOK)
}
