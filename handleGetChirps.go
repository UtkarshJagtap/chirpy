package main

import (
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
)

type valid struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
}

func (api *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	// Get all chirps first
	dbChirps, err := api.databaseQueries.GetChrips(r.Context())
	if err != nil {
		log.Println("can not get chirps", err)
		respondWithError(w, 500, "An internal error with database while getting chirps")
		return
	}

	// Parse author_id if provided
	authorID := uuid.Nil
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, 400, "Invalid author ID")
			return
		}
	}
	sortDirection := "asc"
	sortDirectionParam := r.URL.Query().Get("sort")
	if sortDirectionParam == "desc" {
		sortDirection = "desc"
	}

	// Build the response, filtering by author if needed
	var chirps []valid
	for _, dbChirp := range dbChirps {
		// If author_id was specified, only include chirps from that author
		if authorID != uuid.Nil && dbChirp.UserID != authorID {
			continue
		}

		// Convert to your valid struct format
		chirps = append(chirps, valid{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}
	sort.Slice(chirps, func(i, j int) bool {
		if sortDirection == "desc" {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		}
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	})
	respondWithJSON(w, 200, chirps)
}

func (api *apiConfig) handleGetAChirp(w http.ResponseWriter, r *http.Request) {

	chirp_id := r.PathValue("chirpID")
	id, err := uuid.Parse(chirp_id)

	if err != nil {
		log.Println("there was an error converting id", err)
		respondWithError(w, 404, "can not convert requested chirp id")
		return
	}

	chirp, err := api.databaseQueries.GetChirp(r.Context(), id)
	if err != nil {
		log.Println("There was an error retriving a chirp", err)
		respondWithError(w, 404, "can not retrive requested chirp id")
		return
	}

	actual := valid{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, 200, actual)

}
