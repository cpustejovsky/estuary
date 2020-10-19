package main

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/cpustejovsky/estuary/pkg/models"
)

type Notes = []models.Note

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
		{"Getting Notes", "/api/notes/category/test", "get", http.StatusOK, []byte("Hello")},
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
