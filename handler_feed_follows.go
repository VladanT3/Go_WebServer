package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/VladanT3/Go_WebServer/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
    type parameters struct {
        FeedID uuid.UUID `json:"feed_id"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
        return
    }

    feed_follow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
        ID: uuid.New(),
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        UserID: user.ID,
        FeedID: params.FeedID,
    })
    if err != nil {
        respondWithError(w, 400, fmt.Sprintf("Couldn't follow feed: %v", err))
        return
    }

	respondWithJSON(w, 201, databaseFeedFollowToFeedFollow(feed_follow))
}

func (apiCfg *apiConfig) handlerGetFollowedFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
    followed_feeds, err := apiCfg.DB.GetFeedsByUser(r.Context(), user.ID)

    if err != nil {
    respondWithError(w, 400, fmt.Sprintf("Couldn't get feeds: %v", err))
    return
    }

    respondWithJSON(w, 200, databaseFeedsToFeeds(followed_feeds))
}

func (apiCfg *apiConfig) handlerUnfollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
    var feedID string = chi.URLParam(r, "feedID")
    feedUUID, err := uuid.Parse(feedID)
    if err != nil {
        respondWithError(w, 400, fmt.Sprintf("Error parsing feedID to UUID: %v" ,err))
    }

    err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
        UserID: user.ID,
        FeedID: feedUUID,
    })
    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Error unfollowing feed: %v" ,err))
    }

    respondWithJSON(w, 200, struct{}{})
}
