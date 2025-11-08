package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/utkarshjagtap/chirpy/internal/auth"
	"github.com/utkarshjagtap/chirpy/internal/database"
)

func (api *apiConfig) handleChirpDelete(w http.ResponseWriter, r *http.Request) {

	access_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, "StatusUnauthorized")
		return
	}

	user_id, err := auth.ValidateJWT(access_token, api.jwts)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, "StatusUnauthorized")
		return
	}

	chirp_id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "invalid chirp id")
		return
	}

	chirp, err := api.databaseQueries.GetChirp(r.Context(), chirp_id)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusNotFound, "not found")
		return
	}

	log.Println(chirp.UserID, "and", user_id)
	if chirp.UserID != user_id {
		log.Println(err)
		respondWithError(w, http.StatusForbidden, "StatusForbidden")
		return
	}

	dberr := api.databaseQueries.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirp_id,
		UserID: user_id,
	})

	if dberr != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "can't delete")
		return
	}

	w.WriteHeader(204)

}
