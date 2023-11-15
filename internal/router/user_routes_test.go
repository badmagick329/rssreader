package router

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/badmagick329/rssreader/internal/auth"
	"github.com/badmagick329/rssreader/internal/database"
	"github.com/badmagick329/rssreader/internal/handlers"
)

func TestHealthz(t *testing.T) {
	t.Run("v1/healthz returns 200", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/v1/healthz", nil)
		response := httptest.NewRecorder()
		handlers.HandlerReadiness(response, request)
		got := response.Code
		want := 200
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}

func TestCreateUser(t *testing.T) {
	type Body struct {
		Name string `json:"name"`
	}
	body1, _ := json.Marshal(Body{"tim"})
	requestBody1 := bytes.NewReader(body1)
	emptyBody := bytes.NewReader([]byte{})
	tests := map[string]struct {
		body      io.Reader
		wantCode  int
		wantBody  string
		validName string
	}{
		"valid user": {
			body:      requestBody1,
			wantCode:  http.StatusCreated,
			wantBody:  "",
			validName: "tim",
		},
		"empty body": {
			body:     emptyBody,
			wantCode: http.StatusBadRequest,
			wantBody: "",
		},
		"empty name": {
			body:     bytes.NewReader([]byte(`{"name":""}`)),
			wantCode: http.StatusBadRequest,
			wantBody: "",
		},
	}
	cfg := handlers.New(TEST_DB_URL)
	cfg.ClearDB(context.Background())
	defer cfg.Close()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodPost, "/v1/users", tc.body)
			response := httptest.NewRecorder()
			cfg.HandlerCreateUser(response, request)
			got := response.Code
			if got != tc.wantCode {
				t.Errorf("got %d, want %d", got, tc.wantCode)
			}
			if tc.wantCode == 201 {
				gotUser := handlers.User{}
				json.Unmarshal([]byte(response.Body.String()), &gotUser)
				if gotUser.Name != tc.validName {
					t.Errorf("got %s, want %s", gotUser.Name, tc.validName)
				}
				if len(gotUser.Apikey) != APIKEY_LEN {
					t.Errorf("got apikey len %d, want %d", len(gotUser.Apikey), APIKEY_LEN)
				}
				if len(gotUser.ID) != ID_LEN {
					t.Errorf("got id len %d, want %d", len(gotUser.ID), ID_LEN)
				}

			} else {
				gotBody := response.Body.String()
				if tc.wantBody != "" && gotBody != tc.wantBody {
					t.Errorf("got %s, want %s", gotBody, tc.wantBody)
				}
			}
		})
	}
}

func TestGetUserSuccess(t *testing.T) {
	type TestCase struct {
		apikey   string
		wantCode int
	}
	tests := map[string]TestCase{}
	ctx := context.Background()
	cfg := handlers.New(TEST_DB_URL)
	cfg.ClearDB(ctx)
	users, err := cfg.DB.GetUsers(ctx)
	if err != nil {
		t.Errorf("Could not get users")
	}
	if len(users) != 0 {
		t.Errorf("Expected 0 users, got %d", len(users))
	}
	users = CreateDummyUsers(&cfg, ctx, 3)
	for i, user := range users {
		tests[fmt.Sprintf("TestGetUser_%d", i)] = TestCase{
			apikey:   user.Apikey,
			wantCode: 200,
		}
	}
	defer cfg.Close()
	for name, tc := range tests {
		numString := strings.Split(name, "_")[1]
		num, _ := strconv.Atoi(numString)
		expectedUser := users[num]
		t.Run(name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/v1/users", nil)
			request.Header.Set(auth.API_KEY_HEADER, "ApiKey "+tc.apikey)
			response := httptest.NewRecorder()
			// cfg.HandlerGetUser(response, request)
			cfg.MiddlewareAuth(cfg.HandlerGetUserAuthed)(response, request)
			got := response.Code
			if got != tc.wantCode {
				t.Errorf("got %d, want %d", got, tc.wantCode)
			}
			handlersUser := handlers.User{}
			json.Unmarshal([]byte(response.Body.String()), &handlersUser)
			if handlersUser.Name != expectedUser.Name {
				t.Errorf("got %s, want %s", handlersUser.Name, expectedUser.Name)
			}
			if handlersUser.ID.String() != expectedUser.ID.String() {
				t.Errorf("got %s, want %s", handlersUser.ID, expectedUser.ID.String())
			}
			if handlersUser.Apikey != expectedUser.Apikey {
				t.Errorf("got %s, want %s", handlersUser.Apikey, expectedUser.Apikey)
			}
			cleanedExpectedCreatedAt := strings.Split(expectedUser.CreatedAt.String(), ".")[0]
			cleanedExpectedUpdatedAt := strings.Split(expectedUser.UpdatedAt.String(), ".")[0]
			cleanedGotCreatedAt := strings.Split(handlersUser.CreatedAt.String(), ".")[0]
			cleanedGotUpdatedAt := strings.Split(handlersUser.UpdatedAt.String(), ".")[0]
			if cleanedGotCreatedAt != cleanedExpectedCreatedAt {
				t.Errorf("got %s, want %s", handlersUser.CreatedAt, expectedUser.CreatedAt)
			}
			if cleanedGotUpdatedAt != cleanedExpectedUpdatedAt {
				t.Errorf("got %s, want %s", handlersUser.UpdatedAt, expectedUser.UpdatedAt)
			}
		})
	}
}

func TestGetUserFail(t *testing.T) {
	tests := map[string]struct {
		apikey   string
		wantCode int
	}{
		"no apikey": {
			apikey:   "",
			wantCode: http.StatusBadRequest,
		},
		"invalid apikey": {
			apikey:   "ApiKey 1234567890123456789012345678901234567890123456789012345678901234",
			wantCode: http.StatusUnauthorized,
		},
	}
	cfg := handlers.New(TEST_DB_URL)
	ctx := context.Background()
	cfg.ClearDB(ctx)
	CreateDummyUsers(&cfg, ctx, 3)
	defer cfg.Close()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/v1/users", nil)
			request.Header.Set(auth.API_KEY_HEADER, tc.apikey)
			response := httptest.NewRecorder()
			cfg.MiddlewareAuth(cfg.HandlerGetUserAuthed)(response, request)
			got := response.Code
			if got != tc.wantCode {
				t.Errorf("got %d, want %d", got, tc.wantCode)
			}
		})
	}
}

func CreateDummyUsers(cfg *handlers.Config, ctx context.Context, n int) []database.User {
	users := []database.User{}
	for i := 0; i < n; i++ {
		createParams := handlers.GetUserCreateParams(fmt.Sprintf("user%d", i))
		user, err := cfg.DB.CreateUser(ctx, createParams)
		if err != nil {
			log.Fatalf("Error creating user: %s", err)
		}
		users = append(users, user)
	}
	return users
}
