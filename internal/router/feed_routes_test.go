package router

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/badmagick329/rssreader/internal/auth"
	"github.com/badmagick329/rssreader/internal/handlers"
	"github.com/google/uuid"
)

type FeedResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
}

var testBody = bytes.NewReader([]byte(`{"name":"test","url":"https://www.test.com"}`))

func TestCreateFeed(t *testing.T) {
	cfg := handlers.New(DB_URL)
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
				var resp FeedResponse
				json.Unmarshal([]byte(response.Body.String()), &resp)
				if resp.Name != tc.wantName {
					t.Errorf("got %s, want %s", resp.Name, tc.wantName)
				}
				if resp.Url != tc.wantUrl {
					t.Errorf("got %s, want %s", resp.Url, tc.wantUrl)
				}
				if resp.UserID != tc.wantUserID {
					t.Errorf("got %s, want %s", resp.UserID, tc.wantUserID)
				}
			}
		})
	}
}

