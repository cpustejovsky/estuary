package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestApplication(t *testing.T) *application {
	return &application{
		errorLog: log.New(ioutil.Discard, "", 0),
		infoLog:  log.New(ioutil.Discard, "", 0),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)
	return &testServer{ts}
}

func (ts *testServer) method(t *testing.T, urlPath string, method string) (int, http.Header, []byte) {
	url := ts.URL + urlPath
	reqBody := strings.NewReader(`
	{
		"foo":"bar"
	}`)
	var response *http.Response
	switch method {
	case "get":
		rs, err := ts.Client().Get(url)
		if err != nil {
			t.Fatal(err)
		}
		response = rs
	case "post":
		rs, err := ts.Client().Post(url, "application/json", reqBody)
		if err != nil {
			t.Fatal(err)
		}
		response = rs
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	return response.StatusCode, response.Header, body
}
