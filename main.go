package main

import (
	"log"
	"net/http"
)

func middlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	const PORT = "8080"

	mux := http.NewServeMux()
	corsMux := middlewareCORS(mux)

	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: corsMux,
	}

	log.Printf("Serving on port: %s\n", PORT)
	log.Fatal(server.ListenAndServe())
}
