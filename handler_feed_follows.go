package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	// "github.com/harshit18-09/RSS-Aggregator/internal/auth"
	"github.com/harshit18-09/RSS-Aggregator/internal/db"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "JSON was not parsed")
		return
	}

	feed, err := apiCfg.DB.CreateFeedFollow(r.Context(), db.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    params.FeedID,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, 500, "Failed to create feed follow")
		return
	}

	respondWithJSON(w, 201, databaseFeedFollowToFeedFollow(feed))
}
