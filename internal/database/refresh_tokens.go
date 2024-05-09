package database

import "time"

type RefreshToken struct {
	Token     string    `json:"token"`
	UserID    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (db *DB) CreateRefreshToken(token string, userId int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	refreshToken := RefreshToken{
		Token:     token,
		UserID:    userId,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	dbStructure.RefreshTokens[token] = refreshToken

	return db.writeDB(dbStructure)
}

func (db *DB) GetRefreshToken(token string) (RefreshToken, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	if refreshToken, ok := dbStructure.RefreshTokens[token]; ok {
		return refreshToken, nil
	}

	return RefreshToken{}, ErrNotExist
}

func (db *DB) DeleteRefreshToken(token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(dbStructure.RefreshTokens, token)

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
