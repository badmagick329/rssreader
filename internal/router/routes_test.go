package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	t.Run("/v1/healthz returns 200", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/v1/healthz", nil)
		response := httptest.NewRecorder()
		HandlerReadiness(response, request)
		got := response.Code
		want := 200
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}
