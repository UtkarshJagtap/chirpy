package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/utkarshjagtap/chirpy/internal/auth"
	"github.com/utkarshjagtap/chirpy/internal/database"
)

func (api *apiConfig) handleChirpyRedUpgrade(w http.ResponseWriter, r *http.Request) {

	apikey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Println(err)
		w.WriteHeader(401)
		return
	}

	if apikey != api.polka {
		w.WriteHeader(401)
		return
	}

	type data struct {
		UserId string `json:"user_id"`
	}

	type reqestshape struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	newdecoder := json.NewDecoder(r.Body)
	reqestdata := reqestshape{}
	err = newdecoder.Decode(&reqestdata)

	if err != nil {
		log.Println(err)
		w.WriteHeader(204)
		return
	}

	if reqestdata.Event != "user.upgraded" {
		log.Println(reqestdata.Event)
		w.WriteHeader(204)
		return
	}

	user_id, err := uuid.Parse(reqestdata.Data.UserId)
	if err != nil {
		log.Println(reqestdata.Event)
		w.WriteHeader(204)
		return
	}

	_, err = api.databaseQueries.UpgradeChirpy(r.Context(), database.UpgradeChirpyParams{
		ChirpyRed: true,
		ID:        user_id,
	})

	if err != nil {
		log.Println(err)
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(204)
}
