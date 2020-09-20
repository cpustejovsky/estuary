package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHome(t *testing.T) {
	t.Run("returns 'Hello, World' for '/' route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		Server(response, request)

		got := response.Body.String()
		want := "Hello, World"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestNoteRoutesExist(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes)
	tests := []struct {
		name     string
		urlPath  string
		method   string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/api/notes/category", "get", http.StatusOK, nil},
		{"Non-existent ID", "/api/notes", "post", http.StatusOK, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.method(t, tt.urlPath, tt.method)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}
		})
	}
}
