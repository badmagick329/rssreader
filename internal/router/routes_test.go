package router

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/badmagick329/rssreader/internal/handlers"
)

const DB_URL = "postgres://test:test@localhost:5433/test?sslmode=disable"

func TestHealthz(t *testing.T) {
	t.Run("/v1/healthz returns 200", func(t *testing.T) {
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
		body     io.Reader
		wantCode int
		wantBody string
	}{
		"valid user": {
			body:     requestBody1,
			wantCode: 200,
			wantBody: string(body1),
		},
		"empty body": {
			body:     emptyBody,
			wantCode: 400,
			wantBody: "",
		},
		"empty name": {
			body:     bytes.NewReader([]byte(`{"name":""}`)),
			wantCode: 400,
			wantBody: "",
		},
	}
	cfg := handlers.New(DB_URL)
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
			gotBody := response.Body.String()
			if tc.wantBody != "" && gotBody != tc.wantBody {
				t.Errorf("got %s, want %s", gotBody, tc.wantBody)
			}
		})
	}
}
