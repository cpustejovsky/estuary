package main

import (
	"bytes"
	"net/http"
	"testing"
)

func TestNoteRoutesExist(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	tests := []struct {
		name     string
		urlPath  string
		method   string
		wantCode int
		wantBody []byte
	}{
		{"Home", "/", "get", http.StatusOK, []byte("Hello, World")},
		{"Getting Notes", "/api/notes", "get", http.StatusOK, nil},
		{"Posting to Notes", "/api/notes", "post", http.StatusOK, nil},
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
