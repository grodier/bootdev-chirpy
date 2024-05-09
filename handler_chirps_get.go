package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
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

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:   dbChirp.ID,
		Body: dbChirp.Body,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {

	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	authorId := -1
	authorIdStr := r.URL.Query().Get("author_id")
	if authorIdStr != "" {
		authorId, err = strconv.Atoi(authorIdStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID")
		}
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		if authorId != -1 && authorId != dbChirp.AuthorId {
			continue
		}
		chirps = append(chirps, Chirp{
			ID:       dbChirp.ID,
			Body:     dbChirp.Body,
			AuthorId: dbChirp.AuthorId,
		})
	}

	asc := true
	sortQ := r.URL.Query().Get("sort")
	if sortQ == "desc" {
		asc = false
	}
	sort.Slice(chirps, func(i, j int) bool {
		if asc {
			return chirps[i].ID < chirps[j].ID
		}
		return chirps[j].ID < chirps[i].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}
