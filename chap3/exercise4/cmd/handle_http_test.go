package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func startTestHttpServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/download", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "this is a response")
	})
	mux.HandleFunc("/upload", func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		data, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "JSON request received: %d bytes", len(data))
	})
	return httptest.NewServer(mux)
}

func TestHandleHttp(t *testing.T) {
	usageMessage, err := os.ReadFile("../testdata/expectedGolden.cmd.httpCmdUsage")
	if err != nil {
		t.Fatalf("error reading expectedGolden.cmd.httpCmdUsage")
	}

	ts := startTestHttpServer()
	defer ts.Close()

	outputFile := filepath.Join(t.TempDir(), "file_path.out")
	jsonBody := `{"id":1}`
	jsonBodyFile := filepath.Join(t.TempDir(), "data.json")
	err = os.WriteFile(jsonBodyFile, []byte(jsonBody), 0666)
	if err != nil {
		t.Fatal(err)
	}
	uploadData := "This is some data"
	uploadFile := filepath.Join(t.TempDir(), "file.data")
	err = os.WriteFile(uploadFile, []byte(uploadData), 0666)
	if err != nil {
		t.Fatal(err)
	}

	testConfigs := []struct {
		args   []string
		output string
		err    error
	}{
		{
			args: []string{},
			err:  ErrNoServerSpecified,
		},
		{
			args:   []string{"-h"},
			err:    errors.New("flag: help requested"),
			output: string(usageMessage),
		},
		{
			args:   []string{ts.URL + "/download"},
			err:    nil,
			output: "this is a response\n",
		},
		{
			args:   []string{"-verb", "PUT", "http://localhost"},
			err:    ErrInvalidHTTPMethod,
			output: "Invalid HTTP method\n",
		},
		{
			args:   []string{"-verb", "GET", "-output", outputFile, ts.URL + "/download"},
			err:    nil,
			output: fmt.Sprintf("Data saved to: %s\n", outputFile),
		},
		{
			args:   []string{"-verb", "POST", "-body", "", ts.URL + "/upload"},
			err:    ErrInvalidHTTPPostRequest,
			output: "HTTP POST request must specify a non-empty JSON body\n",
		},
		{
			args:   []string{"-verb", "POST", "-body", jsonBody, ts.URL + "/upload"},
			err:    nil,
			output: fmt.Sprintf("JSON request received: %d bytes\n", len(jsonBody)),
		},
		{
			args:   []string{"-verb", "POST", "-body-file", jsonBodyFile, ts.URL + "/upload"},
			err:    nil,
			output: fmt.Sprintf("JSON request received: %d bytes\n", len(jsonBody)),
		},
		{
			args: []string{
				"-verb", "POST",
				"-upload", uploadFile,
				"-form-data", "filename=test.data",
				"-form-data", "version=0.1",
				ts.URL + "/upload",
			},
			err: nil,
			output: fmt.Sprintf(
				"HTTP POST request received:filename=test.data,version=0.1,upload=%d bytes",
				len(uploadFile),
			),
		},
		{
			args: []string{
				"-verb", "POST",
				"-body-file", jsonBody,
				"-upload", uploadFile,
				"-form-data", "filename=test.data",
				"-form-data", "version=0.1",
				ts.URL + "/upload",
			},
			err: nil,
			output: fmt.Sprintf(
				"HTTP POST request received:json=%d bytes,filename=test.data,version=0.1,upload=%d bytes",
				len(jsonBody), len(uploadFile),
			),
		},
		{
			args: []string{
				"-verb", "POST",
				"-body", jsonBody,
				"-upload", uploadFile,
				"-form-data", "filename=test.data",
				"-form-data", "version=0.1",
				ts.URL + "/upload",
			},
			err: nil,
			output: fmt.Sprintf(
				"HTTP POST request received:json=%d bytes,filename=test.data,version=0.1,upload=%d bytes",
				len(jsonBody), len(uploadFile),
			),
		},
	}
	byteBuf := new(bytes.Buffer)
	for i, tc := range testConfigs {
		t.Log(i)
		err := HandleHttp(byteBuf, tc.args)
		if tc.err == nil && err != nil {
			t.Fatalf("Expected nil error, got %v", err)
		}

		if tc.err != nil && err == nil {
			t.Fatal("Expected non-nil error, got nil")
		}

		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("Expected error %v, got %v", tc.err, err)
		}

		gotOutput := byteBuf.String()
		if tc.output != gotOutput {
			gotOutputFilename := fmt.Sprintf("testdata/gotOutput.cmd.%d", i)
			t.Errorf(
				"Expected output to be:\n%s\n\nGot:\n%s\n\n"+
					"Writing expected data to file: %s",
				tc.output, gotOutput,
				gotOutputFilename,
			)
			if ok := os.WriteFile(gotOutputFilename, []byte(gotOutput), 0666); ok != nil {
				t.Fatal("Error writing expected output to file", err)
			}
		}
		byteBuf.Reset()
	}
}
