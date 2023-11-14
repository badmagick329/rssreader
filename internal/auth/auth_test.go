package auth

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestGetApiKey(t *testing.T) {
	header1 := http.Header{}
	key1 := "1234567890123456789012345678901234567890123456789012345678901234"
	header1.Add(API_KEY_HEADER, fmt.Sprintf("ApiKey %s", key1))
	tests := map[string]struct {
		header http.Header
		want   string
		err    error
	}{
		"valid key": {
			header: header1,
			want:   key1,
			err:    nil,
		},
		"no key": {
			header: http.Header{},
			want:   "",
			err:    errors.New("s"),
		},
		"invalid key": {
			header: http.Header{API_KEY_HEADER: []string{"ApiKey"}},
			want:   "",
			err:    errors.New("s"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := GetApiKey(tc.header)
			if tc.err != nil && err == nil {
				t.Errorf("got nil, want %s", tc.err)
			}
			if tc.err == nil && err != nil {
				t.Errorf("got %s, want nil", err)
			}
			if got != tc.want {
				t.Errorf("got %s, want %s", got, tc.want)
			}
		})
	}

}
