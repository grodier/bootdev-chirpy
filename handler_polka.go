package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerPolka(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusNoContent)
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
