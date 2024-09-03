package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/elishambadi/go-rss-agg/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
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
		UserID:    user.ID,
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

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't get feeds from DB: %v", err))
		return
	}

	respondWithJSON(w, 200, databaseFeedsToFeeds(feeds))

}

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprint("Error passing JSON:", err))
		return
	}

	// DB method defined in sql/queries.sql
	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error creating feed follow: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *apiConfig) handlerGetFeedFollowsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollowsByUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't get feed follows from DB: %v", err))
		return
	}

	respondWithJSON(w, 200, databaseFeedFollowsToFeedFollows(feedFollows))

}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, 400, "Couldn't pass Feed Follow ID")
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't delete feed follow: %v", err))
	}

	respondWithJSON(w, 200, struct{}{})
}

func (apiCfg *apiConfig) handlerDeleteFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	feedIDStr := chi.URLParam(r, "feedId")
	feedID, err := uuid.Parse(feedIDStr)
	if err != nil {
		respondWithError(w, 400, "Couldn't parse Feed ID")
		return
	}

	err = apiCfg.DB.DeleteFeed(r.Context(), database.DeleteFeedParams{
		ID:     feedID,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't delete feed: %v", err))
	}

	respondWithJSON(w, 200, struct{}{})
}
