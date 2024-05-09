package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/grodier/bootdev-chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
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

	refreshToken, err := cfg.DB.GetRefreshToken(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	jwt, err := auth.GenerateJWT(refreshToken.UserID, time.Hour, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Issue generating token")
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: jwt,
	})
}
