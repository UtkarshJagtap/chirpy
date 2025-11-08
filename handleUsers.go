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

func (apicf *apiConfig) handleUsers(w http.ResponseWriter, r *http.Request) {

	type user struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		ChirpyRed bool      `json:"is_chirpy_red"`
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

	// creating the data
	pass, err := auth.HashPassword(reqdata.Pass)
	if err != nil {
		log.Println("hashing error", err)
		respondWithError(w, 500, "Internal error while decoding the request body")
		return

	}

	data, err := apicf.databaseQueries.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          reqdata.Email,
		HashedPassword: pass,
	})

	if err != nil {
		log.Println("57", err)
		respondWithError(w, 500, "Internal Error with database")
		return
	}

	responseUser := user{
		Id:        data.ID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Email:     data.Email,
		ChirpyRed: data.ChirpyRed,
	}

	respondWithJSON(w, 201, responseUser)
}
