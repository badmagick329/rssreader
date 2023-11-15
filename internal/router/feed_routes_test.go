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

var testBody = bytes.NewReader([]byte(`{"name":"test","url":"https://www.test.com"}`))

func TestCreateFeed(t *testing.T) {
	cfg := handlers.New(TEST_DB_URL)
	ctx := context.Background()
	cfg.ClearDB(ctx)
	defer cfg.Close()
	user := CreateDummyUsers(&cfg, ctx, 1)[0]
	tests := map[string]struct {
		body       io.Reader
		headerKey  string
		headerVal  string
		wantCode   int
		wantName   string
		wantUrl    string
		wantUserID uuid.UUID
	}{
		"valid feed": {
			body:       testBody,
			headerKey:  auth.API_KEY_HEADER,
			headerVal:  "ApiKey " + user.Apikey,
			wantCode:   http.StatusCreated,
			wantName:   "test",
			wantUrl:    "https://www.test.com",
			wantUserID: user.ID,
		},
		"No apikey data fails": {
			body:       testBody,
			headerKey:  "",
			headerVal:  "",
			wantCode:   http.StatusBadRequest,
			wantName:   "",
			wantUrl:    "",
			wantUserID: uuid.UUID{},
		},
		"Wrong key fails": {
			body:       testBody,
			headerKey:  auth.API_KEY_HEADER,
			headerVal:  "ApiKey " + "wrongkey",
			wantCode:   http.StatusUnauthorized,
			wantName:   "",
			wantUrl:    "",
			wantUserID: uuid.UUID{},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			request, _ := http.NewRequest("POST", "/api/feeds", tc.body)
			request.Header.Set(tc.headerKey, tc.headerVal)
			response := httptest.NewRecorder()
			cfg.MiddlewareAuth(cfg.HandlerCreateFeed)(response, request)
			got := response.Code
			if got != tc.wantCode {
				t.Errorf("got %d, want %d", got, tc.wantCode)
			}
			if got == http.StatusCreated {
				var gotFeed handlers.Feed
				decoder := json.NewDecoder(response.Body)
				decoder.Decode(&gotFeed)
				if gotFeed.Name != tc.wantName {
					t.Errorf("got %s, want %s", gotFeed.Name, tc.wantName)
				}
				if gotFeed.Url != tc.wantUrl {
					t.Errorf("got %s, want %s", gotFeed.Url, tc.wantUrl)
				}
				if gotFeed.UserID != tc.wantUserID {
					t.Errorf("got %s, want %s", gotFeed.UserID, tc.wantUserID)
				}
			}
		})
	}
}

func TestGetFeeds(t *testing.T) {
	cfg := handlers.New(TEST_DB_URL)
	ctx := context.Background()
	cfg.ClearDB(ctx)
	users := CreateDummyUsers(&cfg, ctx, 1)
	user := users[0]
	// feeds := CreateDummyFeeds(&cfg, ctx, user, 3)
	CreateDummyFeeds(&cfg, ctx, user, 3)
	tests := map[string]struct {
		wantCode int
	}{
		"feeds are returned": {
			wantCode: http.StatusOK,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", "/v1/feeds", nil)
			response := httptest.NewRecorder()
			cfg.HandlerGetFeeds(response, request)
			if response.Code != tc.wantCode {
				t.Errorf("got %d, want %d", response.Code, tc.wantCode)
			}
			var gotFeeds []handlers.Feed
			decoder := json.NewDecoder(response.Body)
			decoder.Decode(&gotFeeds)
			for _, f := range gotFeeds {
				if f.UserID != user.ID {
					t.Errorf("got %s, want %s", f.UserID, user.ID)
				}
			}
		})
	}
}

func CreateDummyFeeds(
	cfg *handlers.Config,
	ctx context.Context,
	user database.User,
	n int,
) []database.Feed {
	feeds := make([]database.Feed, n)
	for i := 0; i < n; i++ {
		feedName := fmt.Sprintf("feed_%d", i)
		feedUrl := fmt.Sprintf("http://%s.com", feedName)
		feedParams := handlers.GetFeedCreateParams(feedName, feedUrl, user.ID)
		feed, _ := cfg.DB.CreateFeed(ctx, feedParams)
		feeds[i] = feed
	}
	return feeds

}
