package database

import "time"

type RefreshToken struct {
	Token  string    `json:"token"`
	UserID int       `json:"user_id"`
	Exp    time.Time `json:"exp"`
}

func (db *DB) CreateRefreshToken(token string, userId int) (RefreshToken, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	refreshToken := RefreshToken{
		Token:  token,
		UserID: userId,
		Exp:    time.Now().Add(time.Hour * 24 * 60),
	}
	dbStructure.RefreshTokens[token] = refreshToken

	err = db.writeDB(dbStructure)
	if err != nil {
		return RefreshToken{}, err
	}

	return refreshToken, nil
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
