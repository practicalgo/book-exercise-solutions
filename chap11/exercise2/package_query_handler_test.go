package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "gocloud.dev/blob/fileblob"
)

func TestPackageQueryHandlerNoResults(t *testing.T) {
	packageBucket, err := getTestBucket(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer packageBucket.Close()

	config := appConfig{
		logger: log.New(
			os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile,
		),
		packageBucket: packageBucket,
		db:            testDb,
	}

	mux := http.NewServeMux()
	setupHandlers(mux, config)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := http.Client{}
	resp, err := client.Get(ts.URL + "/api/packages?owner_id=1&name=test-package&version=0.1")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf(
			"Expected a HTTP 200 response, Got: %v\n",
			resp.StatusCode,
		)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	expectedResponse := `{"packages":[]}`
	if string(data) != expectedResponse {
		t.Fatalf("Expected response to be: %s, Got: %s", expectedResponse, string(data))
	}
}

func TestPackageQueryHandler(t *testing.T) {
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
			OwnerId:       2,
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

	client := http.Client{}
	resp, err := client.Get(ts.URL + "/api/packages?owner_id=1&name=pkg")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf(
			"Expected a HTTP 200 response, Got: %v\n",
			resp.StatusCode,
		)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var gotResponse pkgQueryResponse
	err = json.Unmarshal(data, &gotResponse)
	if err != nil {
		t.Fatal(err)
	}
	expectedResponseJson := `{"packages":[{"owner_id":2,"name":"pkg","version":"0.1","object_store_id":"pkg-0.1-pkg-0.1.tar.gz","created":"2021-11-30 08:38:00"}]}`
	var expectedResponse pkgQueryResponse
	err = json.Unmarshal([]byte(expectedResponseJson), &expectedResponse)
	if err != nil {
		t.Fatal(err)
	}

	if len(gotResponse.Packages) != len(expectedResponse.Packages) {
		t.Fatalf("Expected response to contain %d packages, Got: %d", len(expectedResponse.Packages), len(gotResponse.Packages))
	}
}
