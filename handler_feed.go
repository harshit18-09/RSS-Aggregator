package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	// "github.com/harshit18-09/RSS-Aggregator/internal/auth"
	"github.com/harshit18-09/RSS-Aggregator/internal/db"
)

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "JSON was not parsed")
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), db.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, 500, "Failed to create feed")
		return
	}

	respondWithJSON(w, 201, databaseFeedToFeed(feed))
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeed(r.Context())
	if err != nil {
		respondWithError(w, 500, "Failed to get feeds")
		return
	}
	respondWithJSON(w, 201, databaseFeedsToFeeds(feeds))
}
