package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	_ "gocloud.dev/blob/fileblob"
)

func TestPackageGetHandler(t *testing.T) {
	packageBucket, err := getTestBucket(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer packageBucket.Close()

	// create a test object
	err = packageBucket.WriteAll(
		context.Background(),
		"test-object-id",
		[]byte("test-data"),
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	config := appConfig{
		logger: log.New(
			os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile,
		),
		packageBucket: packageBucket,
	}

	mux := http.NewServeMux()
	setupHandlers(mux, config)

	ts := httptest.NewServer(mux)
	defer ts.Close()
	var redirectUrl string
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			redirectUrl = req.URL.String()
			return errors.New("no redirect")
		},
	}

	_, err = client.Get(ts.URL + "/api/packages?id=test-object-id")
	if err == nil {
		t.Fatal("Expected error: no redirect, Got nil")
	}
	if !strings.HasPrefix(redirectUrl, "file:///") {
		t.Fatalf("Expected redirect url to start with file:///, got: %v", redirectUrl)
	}
}

func TestPackageDownload(t *testing.T) {
	packageBucket, err := getTestBucket(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer packageBucket.Close()

	testContents := []byte("test-data")

	// create a test object
	err = packageBucket.WriteAll(
		context.Background(),
		"test-object-id",
		testContents,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	config := appConfig{
		logger: log.New(
			os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile,
		),
		packageBucket: packageBucket,
	}

	mux := http.NewServeMux()
	setupHandlers(mux, config)

	ts := httptest.NewServer(mux)
	defer ts.Close()
	client := http.Client{}

	resp, err := client.Get(ts.URL + "/api/packages?id=test-object-id&download=true")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	contentTypeHeaderValue := resp.Header.Get("Content-Type")
	if !strings.Contains(contentTypeHeaderValue, "text/plain") {
		t.Fatalf("Expected Content-Type response header to be set to text/plain, got: %s", contentTypeHeaderValue)
	}

	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Disposition
	contentDispositionHeaderValue := resp.Header.Get("Content-Disposition")
	if contentDispositionHeaderValue != "attachment; filename=test-object-id" {
		t.Fatalf("Expected Content-Disposition header to be set to: attachment; filename=test-object-id, got: %s", contentDispositionHeaderValue)
	}

	if !bytes.Equal(data, testContents) {
		t.Fatalf("File contents do not match")
	}
}
