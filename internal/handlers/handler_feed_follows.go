package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/badmagick329/rssreader/internal/database"
	"github.com/badmagick329/rssreader/internal/utils"
	"github.com/google/uuid"
)

func (cfg *Config) HandlerCreateFeedFollow(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	savedFeed, err := GetSavedFeed(cfg, r.Context(), r.Body)
	if err != nil {
		utils.RespondWithError(w, 400, err.Error())
		return
	}
	createParams := GetFeedFollowCreateParams(savedFeed.ID, user.ID)
	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), createParams)
	if err != nil {
		utils.RespondWithError(w, 500, "Error creating feed follow")
		return
	}
	feedFollowResponse := databaseFeedFollowToFeedFollow(feedFollow)
	utils.RespondWithJSON(w, 201, feedFollowResponse)
}

func (cfg *Config) HandlerDeleteFeedFollow(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	savedFeed, err := GetSavedFeed(cfg, r.Context(), r.Body)
	if err != nil {
		utils.RespondWithError(w, 400, err.Error())
		return
	}
	feedFollow, err := cfg.DB.GetFeedFollowByFeedID(r.Context(), savedFeed.ID)
	if err != nil {
		utils.RespondWithError(w, 400, "Feed follow does not exist")
		return
	}
	deletedFeedFollow, err := cfg.DB.DeleteFeedFollow(r.Context(), feedFollow.ID)
	if err != nil {
		utils.RespondWithError(w, 500, "Error deleting feed follow")
		return
	}
	feedFollowResponse := databaseFeedFollowToFeedFollow(deletedFeedFollow)
	utils.RespondWithJSON(w, 200, feedFollowResponse)
}

func (cfg *Config) HandlerGetUserFeedFollows(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	dbFeeds, err := cfg.DB.GetFeedsByUserID(r.Context(), user.ID)
	if err != nil {
		utils.RespondWithError(w, 500, "Error getting feeds")
		return
	}
	feeds := DatabaseFeedsToFeeds(dbFeeds)
	utils.RespondWithJSON(w, 200, feeds)
}

func GetSavedFeed(
	cfg *Config,
	ctx context.Context,
	body io.ReadCloser,
) (database.Feed, error) {
	type parameters struct {
		FeedID string `json:"feed_id"`
	}
	decoder := json.NewDecoder(body)
	defer body.Close()
	var params parameters
	decoder.Decode(&params)
	if params == (parameters{}) {
		return database.Feed{}, errors.New("Empty body")
	}
	feedID, err := uuid.Parse(params.FeedID)
	if err != nil {
		return database.Feed{}, errors.New("Invalid feed_id")
	}
	savedFeed, err := cfg.DB.GetFeedByID(ctx, feedID)
	if err != nil {
		return database.Feed{}, errors.New("Feed does not exist")
	}
	return savedFeed, nil
}
