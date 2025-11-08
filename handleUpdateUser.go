package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/utkarshjagtap/chirpy/internal/auth"
	"github.com/utkarshjagtap/chirpy/internal/database"
)

func (api *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type validuser struct {
		ID        uuid.UUID `json:"id"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	access_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user_id, err := auth.ValidateJWT(access_token, api.jwts)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	reqbod := body{}
	err = decoder.Decode(&reqbod)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "can not decode body")
		return
	}

	hashed_pass, err := auth.HashPassword(reqbod.Password)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "can not decode body")
		return
	}

	user, err := api.databaseQueries.UpdatePass(r.Context(), database.UpdatePassParams{
		Email:          reqbod.Email,
		HashedPassword: hashed_pass,
		UpdatedAt:      time.Now(),
		ID:             user_id,
	})

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "can not update")
		return
	}

	response := validuser{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, http.StatusOK, response)

}
