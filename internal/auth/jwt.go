package auth

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userId int, expiresInSeconds *int, secret string) (string, error) {
	expires := 60 * 60 * 24 // one day sec * min * hour
	if expiresInSeconds != nil {
		if *expiresInSeconds < expires {
			expires = *expiresInSeconds
		}
	}
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(expires))),
		Subject:   strconv.Itoa(userId),
	})

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GenerateRefreshTokenString() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Reader.Read(bytes)
	if err != nil {
		return "", err
	}

	refreshTokenString := hex.EncodeToString(bytes)
	return refreshTokenString, nil

}
