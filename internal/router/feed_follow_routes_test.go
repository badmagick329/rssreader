package router

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/badmagick329/rssreader/internal/auth"
	"github.com/badmagick329/rssreader/internal/database"
	"github.com/badmagick329/rssreader/internal/handlers"
	"github.com/google/uuid"
)

func TestCreateFeedFollow(t *testing.T) {
	cfg := handlers.New(TEST_DB_URL)
	ctx := context.Background()
	cfg.ClearDB(ctx)
	defer cfg.Close()
	user := CreateDummyUsers(&cfg, ctx, 1)[0]
	feed := CreateDummyFeeds(&cfg, ctx, user, 1)[0]
	tests := map[string]struct {
		body       io.Reader
		headerKey  string
		headerVal  string
		wantCode   int
		wantFeedID uuid.UUID
		wantUserID uuid.UUID
	}{
		"valid feed follow": {
			body:       bytes.NewReader([]byte(fmt.Sprintf(`{"feed_id":"%s"}`, feed.ID))),
			headerKey:  auth.API_KEY_HEADER,
			headerVal:  "ApiKey " + user.Apikey,
			wantCode:   http.StatusCreated,
			wantFeedID: feed.ID,
			wantUserID: user.ID,
		},
		"invalid feed id fails": {
			body: bytes.NewReader(
				[]byte(fmt.Sprintf(`{"feed_id":"%s"}`, uuid.New().String())),
			),
			headerKey:  auth.API_KEY_HEADER,
			headerVal:  "ApiKey " + user.Apikey,
			wantCode:   http.StatusBadRequest,
			wantFeedID: uuid.UUID{},
			wantUserID: uuid.UUID{},
		},
		"no feed id fails": {
			body:       bytes.NewReader([]byte(`{"feed_id":""}`)),
			headerKey:  auth.API_KEY_HEADER,
			headerVal:  "ApiKey " + user.Apikey,
			wantCode:   http.StatusBadRequest,
			wantFeedID: uuid.UUID{},
			wantUserID: uuid.UUID{},
		},
		"no apikey data fails": {
			body:       bytes.NewReader([]byte(fmt.Sprintf(`{"feed_id":"%s"}`, feed.ID))),
			headerKey:  "",
			headerVal:  "",
			wantCode:   http.StatusBadRequest,
			wantFeedID: uuid.UUID{},
			wantUserID: uuid.UUID{},
		},
		"wrong key fails": {
			body:       bytes.NewReader([]byte(fmt.Sprintf(`{"feed_id":"%s"}`, feed.ID))),
			headerKey:  auth.API_KEY_HEADER,
			headerVal:  "ApiKey " + "wrongkey",
			wantCode:   http.StatusUnauthorized,
			wantFeedID: uuid.UUID{},
			wantUserID: uuid.UUID{},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			request, _ := http.NewRequest("POST", "/v1/feed_follows", tc.body)
			request.Header.Set(tc.headerKey, tc.headerVal)
			response := httptest.NewRecorder()
			cfg.MiddlewareAuth(cfg.HandlerCreateFeedFollow)(response, request)
			got := response.Code
			if got != tc.wantCode {
				t.Errorf("got %d, want %d", got, tc.wantCode)
			}
			if got == http.StatusCreated {
				var gotFeedFollow handlers.FeedFollow
				decoder := json.NewDecoder(response.Body)
				decoder.Decode(&gotFeedFollow)
				if gotFeedFollow.FeedID != tc.wantFeedID {
					t.Errorf("got %s, want %s", gotFeedFollow.FeedID, tc.wantFeedID)
				}
				if gotFeedFollow.UserID != tc.wantUserID {
					t.Errorf("got %s, want %s", gotFeedFollow.UserID, tc.wantUserID)
				}
			}
		})
	}
}

func TestDeleteFeedFollow(t *testing.T) {
	cfg := handlers.New(TEST_DB_URL)
	ctx := context.Background()
	cfg.ClearDB(ctx)
	defer cfg.Close()
	dbUser := CreateDummyUsers(&cfg, ctx, 1)[0]
	user := handlers.DatabaseUserToUser(dbUser)
	dbFeed := CreateDummyFeeds(&cfg, ctx, dbUser, 1)[0]
	feed := handlers.DatabaseFeedToFeed(dbFeed)
	GetDummyFeedFollows(&cfg, ctx, user, feed, 1)
	tests := map[string]struct {
		body       io.Reader
		headerKey  string
		headerVal  string
		wantCode   int
		wantFeedID uuid.UUID
		wantUserID uuid.UUID
	}{
		"valid feed follow delete": {
			body:       bytes.NewReader([]byte(fmt.Sprintf(`{"feed_id":"%s"}`, feed.ID))),
			headerKey:  auth.API_KEY_HEADER,
			headerVal:  "ApiKey " + user.Apikey,
			wantCode:   http.StatusOK,
			wantFeedID: feed.ID,
			wantUserID: user.ID,
		},
		"invalid feed id fails": {
			body: bytes.NewReader(
				[]byte(fmt.Sprintf(`{"feed_id":"%s"}`, uuid.New().String())),
			),
			headerKey:  auth.API_KEY_HEADER,
			headerVal:  "ApiKey " + user.Apikey,
			wantCode:   http.StatusBadRequest,
			wantFeedID: uuid.UUID{},
			wantUserID: uuid.UUID{},
		},
		"no feed id fails": {
			body:       bytes.NewReader([]byte(`{"feed_id":""}`)),
			headerKey:  auth.API_KEY_HEADER,
			headerVal:  "ApiKey " + user.Apikey,
			wantCode:   http.StatusBadRequest,
			wantFeedID: uuid.UUID{},
			wantUserID: uuid.UUID{},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			request, _ := http.NewRequest("DELETE", "/v1/feed_follows", tc.body)
			request.Header.Set(tc.headerKey, tc.headerVal)
			response := httptest.NewRecorder()
			cfg.MiddlewareAuth(cfg.HandlerDeleteFeedFollow)(response, request)
			got := response.Code
			if got != tc.wantCode {
				t.Errorf("got %d, want %d", got, tc.wantCode)
			}
			if got == http.StatusOK {
				var gotFeedFollow handlers.FeedFollow
				decoder := json.NewDecoder(response.Body)
				decoder.Decode(&gotFeedFollow)
				if gotFeedFollow.FeedID != tc.wantFeedID {
					t.Errorf("got %s, want %s", gotFeedFollow.FeedID, tc.wantFeedID)
				}
				if gotFeedFollow.UserID != tc.wantUserID {
					t.Errorf("got %s, want %s", gotFeedFollow.UserID, tc.wantUserID)
				}
			}
		})
	}
}

func GetDummyFeedFollows(
	cfg *handlers.Config,
	ctx context.Context,
	user handlers.User,
	feed handlers.Feed,
	count int,
) []database.FeedFollow {
	feedFollows := make([]database.FeedFollow, count)
	for i := 0; i < count; i++ {
		feedFollow, _ := cfg.DB.CreateFeedFollow(
			ctx,
			handlers.GetFeedFollowCreateParams(feed.ID, user.ID),
		)
		feedFollows[i] = feedFollow
	}
	return feedFollows
}
