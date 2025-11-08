package main

import (
	"log"
	"net/http"
	"time"

	"github.com/utkarshjagtap/chirpy/internal/auth"
)

func (api *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	type valid struct {
		Token string `json:"token"`
	}
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user_id, err := api.databaseQueries.GetUserFromRefreshToken(r.Context(), refresh_token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not found")
		return
	}

	token, err := auth.MakeJWT(
		user_id,
		api.jwts,
		time.Hour,
	)

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "can not create token")
	}

	respondWithJSON(w, http.StatusOK, valid{
		Token: token,
	})

}
