package main

import (
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
