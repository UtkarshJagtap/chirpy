package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/utkarshjagtap/chirpy/internal/database"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	jwts := os.Getenv("JWT")
	polka := os.Getenv("POLKA_KEY")
	log.Println(dbURL)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Println(err)
		return
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		platform:        os.Getenv("PLATFORM"),
		fileserverHits:  atomic.Int32{},
		databaseQueries: dbQueries,
		jwts:            jwts,
		polka:           polka,
	}

	fileshandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	servermux := http.NewServeMux()

	servermux.Handle("/app/", apiCfg.middlewareMatricsInc(fileshandler))
	servermux.HandleFunc("GET /admin/healthz", handleHealth)
	servermux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	servermux.HandleFunc("POST /admin/reset", apiCfg.handleReset)
	servermux.HandleFunc("POST /api/chirps", apiCfg.handleCreateChrip)
	servermux.HandleFunc("GET /api/chirps", apiCfg.handleGetChirps)
	servermux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handleGetAChirp)
	servermux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handleChirpDelete)
	servermux.HandleFunc("POST /api/users", apiCfg.handleUsers)
	servermux.HandleFunc("POST /api/login", apiCfg.handleLogin)
	servermux.HandleFunc("POST /api/refresh", apiCfg.handleRefresh)
	servermux.HandleFunc("POST /api/revoke", apiCfg.handleRevoke)
	servermux.HandleFunc("PUT /api/users", apiCfg.handleUpdateUser)
	servermux.HandleFunc("POST /api/polka/webhooks", apiCfg.handleChirpyRedUpgrade)
	server := http.Server{
		Handler: servermux,
		Addr:    ":8080",
	}

	err = server.ListenAndServe()

	if err != nil {
		fmt.Println(err)
	}

}
