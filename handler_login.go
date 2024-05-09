package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/grodier/bootdev-chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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

	jwt, err := auth.GenerateJWT(user.ID, time.Hour, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Issue generating token")
		return
	}

	refreshToken, err := auth.GenerateRefreshTokenString()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server issue")
		return
	}

	err = cfg.DB.CreateRefreshToken(refreshToken, user.ID)
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
		RefreshToken: refreshToken,
	})
}
