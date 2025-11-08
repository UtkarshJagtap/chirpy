package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/utkarshjagtap/chirpy/internal/auth"
	"github.com/utkarshjagtap/chirpy/internal/database"
)

func (api *apiConfig) handleCreateChrip(w http.ResponseWriter, r *http.Request) {

	type body struct {
		Body string `json:"body"`
	}

	type valid struct {
		Id         uuid.UUID `json:"id"`
		Body       string    `json:"body"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		User_ID    uuid.UUID `json:"user_id"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	userid, err := auth.ValidateJWT(token, api.jwts)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	data := body{}
	err = decoder.Decode(&data)

	// returning response if there is an error while decoding request body
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	//checking if the request body has less than 140 characters

	if len([]rune(data.Body)) > 140 {
		respondWithError(w, 400, "Chrip is too long")
		return
	}

	validr := valid{
		Body:    cleanProfane(data.Body),
		User_ID: userid,
	}

	// valid response
	chrip, err := api.databaseQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      validr.Body,
		UserID:    validr.User_ID,
	})

	if err != nil {
		respondWithError(w, 500, "There was an error creating Chrip")
		log.Println("Error creating Chrip", err)
		return
	}

	validr = valid{
		Id:         chrip.ID,
		Created_at: chrip.CreatedAt,
		Updated_at: chrip.UpdatedAt,
		Body:       chrip.Body,
		User_ID:    chrip.UserID,
	}

	respondWithJSON(w, 201, validr)
	return

}

func cleanProfane(og string) string {
	words := strings.Fields(og)
	for index, word := range words {
		tolower := strings.ToLower(word)
		if tolower == "kerfuffle" || tolower == "sharbert" || tolower == "fornax" {
			words[index] = "****"
		}
	}
	return strings.Join(words, " ")
}

func respondWithError(w http.ResponseWriter, code int, msg string) {

	w.Header().Set("Content-Type", "application/json")

	type errorResponse struct {
		Error string `json:"error"`
	}

	errormsg := errorResponse{
		Error: msg,
	}

	actualresponse, err := json.Marshal(errormsg)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"error":"Internal Error occured"}`))
		return
	}

	w.WriteHeader(code)
	w.Write(actualresponse)

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	jsonresponse, err := json.Marshal(payload)

	if err != nil {
		respondWithError(w, 500, "Internal Error while Marshalling")
		return
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonresponse)

}
