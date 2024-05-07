package main

import (
	"encoding/json"
	"net/http"

	"github.com/grodier/bootdev-chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds *int   `json:"expires_in_seconds,omitempty"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	jwt, err := auth.GenerateJWT(user.ID, params.ExpiresInSeconds, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Issue generating token")
		return
	}

	refreshTokenString, err := auth.GenerateRefreshTokenString()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server issue")
		return
	}

	refreshToken, err := cfg.DB.CreateRefreshToken(refreshTokenString, user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server issue")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
		Token:        jwt,
		RefreshToken: refreshToken.Token,
	})
}
