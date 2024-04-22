package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const FILE_ROOT_PATH = "."
	const PORT = "8080"
	var apiCfg apiConfig

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(FILE_ROOT_PATH))))
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerResetMetrics)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", FILE_ROOT_PATH, PORT)
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
    <html>
      <body>
        <h1>Welcome, Chirpy Admin</h1>
        <p>Chirpy has been visited %d times!</p>
      </body>
    <html>
  `, cfg.fileserverHits)))
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type errorJson struct {
		Error string `json:"error"`
	}

	type successJson struct {
		Valid bool `json:"valid"`
	}

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errBody := errorJson{
			Error: fmt.Sprintf("Invalid parameters: %s", err),
		}
		data, marshErr := json.Marshal(errBody)
		w.WriteHeader(500)
		if marshErr != nil {
			return
		}
		w.Write(data)
		return
	}

	if len(params.Body) > 140 {
		errBody := errorJson{
			Error: "Chirp is too long",
		}
		data, marshErr := json.Marshal(errBody)
		if marshErr != nil {
			w.WriteHeader(500)
		} else {
			w.WriteHeader(400)
			w.Write(data)
		}
		return
	}

	successBody := successJson{
		Valid: true,
	}

	data, marshErr := json.Marshal(successBody)
	if marshErr != nil {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
