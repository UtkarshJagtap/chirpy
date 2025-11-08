package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/utkarshjagtap/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits  atomic.Int32
	databaseQueries *database.Queries
	platform        string
	jwts            string
	polka           string
}

func (api *apiConfig) middlewareMatricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (api *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	if api.platform != "dev" {
		respondWithError(w, 403, "Forbidden")
		return
	}

	api.fileserverHits.Store(0)
	err := api.databaseQueries.DeleteUsers(r.Context())
	if err != nil {
		log.Println(err)
		respondWithError(w, 500, "Internal error while deleting")
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (api *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`,
		api.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(message))
}
