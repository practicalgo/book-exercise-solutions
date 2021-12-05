package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	_ "gocloud.dev/blob/fileblob"
)

func TestPackageDownloadHandler(t *testing.T) {
	packageBucket, err := getTestBucket(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer packageBucket.Close()

	testObjectId := "pkg-0.1-pkg-0.1.tar.gz"

	// create a test object
	err = packageBucket.WriteAll(
		context.Background(),
		testObjectId, []byte("test-data"),
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
		db:            testDb,
	}

	// udpate package metadata for the test object
	err = updateDb(
		config,
		pkgRow{
			OwnerId:       1,
			Name:          "pkg",
			Version:       "0.1",
			ObjectStoreId: testObjectId,
		},
	)
	if err != nil {
		t.Fatal(err)
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

	_, err = client.Get(ts.URL + "/packages/download?owner_id=1&name=pkg&version=0.1")
	if err == nil {
		t.Fatal("Expected error: no redirect, Got nil")
	}
	if !strings.HasPrefix(redirectUrl, "file:///") {
		t.Fatalf("Expected redirect url to start with file:///, got: %v", redirectUrl)
	}
}
