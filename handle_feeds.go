package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/elishambadi/go-rss-agg/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeeds(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name   string    `json:"name"`
		Url    string    `json:"url"`
		UserId uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprint("Error passing JSON:", err))
		return
	}

	// DB method defined in sql/queries.sql
	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    params.UserId,
	})

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error creating feed: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseFeedToFeed(feed))
}

func (apiCfg *apiConfig) handlerGetFeed(w http.ResponseWriter, r *http.Request) {

	// Proceed to adding a user middleware so we can get a user

	// feed, err := apiCfg.DB.GetFeedsByUser(r.Context(), userId)
	// if err != nil {
	// 	respondWithError(w, 400, fmt.Sprintf("error getting user from DB: %v", err))
	// 	return
	// }

	// respondWithJSON(w, 200, databaseFeedToFeed(feed))
}
