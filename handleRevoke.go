package main

import (
	"log"
	"net/http"
	"time"

	"github.com/utkarshjagtap/chirpy/internal/auth"
	"github.com/utkarshjagtap/chirpy/internal/database"
)

func (api *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > 0 {
		log.Println("request has body, can't revoke")
		respondWithError(w, http.StatusBadRequest, "no body requests only")
		return
	}

	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "no body requests only")
		return
	}

	err = api.databaseQueries.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		UpdatedAt: time.Now(),
		Token:     refresh_token,
	})

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "can't delete the refresh token")
		return
	}

	w.WriteHeader(204)

}
