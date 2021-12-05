package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestUserHandlerApiTimeout(t *testing.T) {

	logBuf := new(bytes.Buffer)
	timeoutDuration := 30 * time.Second
	logger := log.New(logBuf, "", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()
	setupHandlers(mux, timeoutDuration, logger)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := http.Client{
		Timeout: 4 * time.Second,
	}
	_, err := client.Get(ts.URL + "/api/users/" + "?ping_server=" + ts.URL)
	if err == nil {
		t.Fatalf("Expected nil error, got:%v", err)
	}
	expectedServerLogLine := "Aborting request processing: context canceled"
	if !strings.Contains(logBuf.String(), expectedServerLogLine) {
		t.Fatalf("Expected server log to contain: %s\n Got: %s", expectedServerLogLine, logBuf.String())
	}

}
