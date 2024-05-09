package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/grodier/bootdev-chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
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

	jwtToken := parts[1]

	subject, err := auth.ValidateJWT(jwtToken, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
	}
	parsedUserId, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid user id")
		return
	}

	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}
	if dbChirp.AuthorId != parsedUserId {
		respondWithError(w, http.StatusForbidden, "Not authorized")
		return
	}

	err = cfg.DB.DeleteChirp(dbChirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp")
		return
	}

	w.WriteHeader(http.StatusOK)
}
