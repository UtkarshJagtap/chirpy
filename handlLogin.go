package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/utkarshjagtap/chirpy/internal/auth"
	"github.com/utkarshjagtap/chirpy/internal/database"
)

func (api *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type user struct {
		Id           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
		ChirpyRed    bool      `json:"is_chirpy_red"`
	}
	// reading the request
	type req struct {
		Email string `json:"email"`
		Pass  string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	reqdata := req{}
	err := decoder.Decode(&reqdata)
	if err != nil {
		respondWithError(w, 500, "Internal error while decoding the request body")
		return
	}

	rawuser, err := api.databaseQueries.GetPass(r.Context(), reqdata.Email)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	err = auth.CheckPasswordHash(reqdata.Pass, rawuser.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	expiration_time := time.Hour

	accessToken, err := auth.MakeJWT(
		rawuser.ID,
		api.jwts,
		expiration_time,
	)

	if err != nil {
		log.Println("Couldn't create access JWT")
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT")
		return
	}

	expires_at := time.Now().AddDate(0, 0, 60)
	reftokenstring, err := auth.MakeRefreshToken()
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "can not create refresh token")
		return
	}
	refresh_token, err := api.databaseQueries.CreateRefreshToken(r.Context(),
		database.CreateRefreshTokenParams{
			Token:     reftokenstring,
			UserID:    rawuser.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			ExpiresAt: expires_at,
			RevokedAt: sql.NullTime{},
		},
	)

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "something to do with databse")
		return
	}

	validuser := user{
		Id:           rawuser.ID,
		CreatedAt:    rawuser.CreatedAt,
		UpdatedAt:    rawuser.UpdatedAt,
		Email:        rawuser.Email,
		Token:        accessToken,
		RefreshToken: refresh_token,
		ChirpyRed:    rawuser.ChirpyRed,
	}

	respondWithJSON(w, 200, validuser)
}
