package main

import (
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header is required")
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header must be in the format 'Bearer {token}'")
		return
	}

	token := parts[1]

	err := cfg.DB.DeleteRefreshToken(token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
