package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/practicalgo/book-exercise-solutions/chap6/exercise3/config"
)

func TestApiHandler(t *testing.T) {
	r := httptest.NewRequest("GET", "/api", nil)
	w := httptest.NewRecorder()

	b := new(bytes.Buffer)
	c := config.InitConfig(b)

	apiHandler(w, r, c)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected response status: %v, Got: %v\n", http.StatusOK, resp.StatusCode)
	}

	expectedResponseBody := "Hello, world!"

	if string(body) != expectedResponseBody {
		t.Errorf("Expected response: %s, Got: %s\n", expectedResponseBody, string(body))
	}
}

func TestHealthcheckHandler(t *testing.T) {

	testConfigs := []struct {
		httpMethod           string
		expectedStatus       int
		expectedResponseBody string
	}{
		{
			httpMethod:           "GET",
			expectedStatus:       http.StatusOK,
			expectedResponseBody: "ok",
		},
		{
			httpMethod:           "POST",
			expectedStatus:       http.StatusMethodNotAllowed,
			expectedResponseBody: "Method not allowed\n",
		},
		{
			httpMethod:           "PUT",
			expectedStatus:       http.StatusMethodNotAllowed,
			expectedResponseBody: "Method not allowed\n",
		},
	}

	for _, tc := range testConfigs {
		r := httptest.NewRequest(tc.httpMethod, "/healthz", nil)
		w := httptest.NewRecorder()

		b := new(bytes.Buffer)
		c := config.InitConfig(b)

		healthCheckHandler(w, r, c)

		resp := w.Result()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error reading response body: %v", err)
		}

		if resp.StatusCode != tc.expectedStatus {
			t.Errorf("Expected response status: %v, Got: %v\n", http.StatusOK, resp.StatusCode)
		}

		if string(body) != tc.expectedResponseBody {
			t.Errorf("Expected response: %s, Got: %s\n", tc.expectedResponseBody, string(body))
		}
	}
}
