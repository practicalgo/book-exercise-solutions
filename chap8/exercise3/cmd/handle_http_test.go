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
	"strings"
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
	mux.HandleFunc("/redirect", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/new-url", http.StatusMovedPermanently)
	})
	mux.HandleFunc("/debug-header-response", func(w http.ResponseWriter, req *http.Request) {
		headers := []string{}
		for k, v := range req.Header {
			if strings.HasPrefix(k, "Debug") {
				headers = append(headers, fmt.Sprintf("%s=%s", k, v[0]))
			}
		}
		fmt.Fprint(w, strings.Join(headers, " "))
	})
	mux.HandleFunc("/debug-basicauth", func(w http.ResponseWriter, req *http.Request) {
		u, p, ok := req.BasicAuth()
		if !ok {
			http.Error(w, "Basic auth missing/malformed", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "%s=%s", u, p)
	})

	return httptest.NewServer(mux)
}

func TestHandleHttp(t *testing.T) {
	usageMessage := `
http: A HTTP client.

http: <options> server

Options: 
  -basicauth string
    	Add basic auth (username:password) credentials to the outgoing request
  -body string
    	JSON data for HTTP POST request
  -body-file string
    	File containing JSON data for HTTP POST request
  -disable-redirect
    	Do not follow redirection request
  -header value
    	Add one or more headers to the outgoing request (key=value)
  -max-idle-conns int
    	Maximum number of idle connections for the connection pool
  -num-requests int
    	Number of requests to make (default 1)
  -output string
    	File path to write the response into
  -verb string
    	HTTP method (default "GET")
`
	ts := startTestHttpServer()
	defer ts.Close()

	outputFile := filepath.Join(t.TempDir(), "file_path.out")
	jsonBody := `{"id":1}`
	jsonBodyFile := filepath.Join(t.TempDir(), "data.json")
	err := os.WriteFile(jsonBodyFile, []byte(jsonBody), 0666)
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
			err:  InvalidInputError{ErrNoServerSpecified},
		},
		{
			args:   []string{"-h"},
			err:    errors.New("flag: help requested"),
			output: usageMessage,
		},
		{
			args:   []string{ts.URL + "/download"},
			err:    nil,
			output: "this is a response\n",
		},
		{
			args: []string{"-verb", "PUT", "http://localhost"},
			err:  InvalidInputError{ErrInvalidHTTPMethod},
		},
		{
			args:   []string{"-verb", "GET", "-output", outputFile, ts.URL + "/download"},
			err:    nil,
			output: fmt.Sprintf("Data saved to: %s\n", outputFile),
		},
		{
			args: []string{"-verb", "POST", "-body", "", ts.URL + "/upload"},
			err:  InvalidInputError{ErrInvalidHTTPPostRequest},
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
			args:   []string{"-verb", "POST", "-body-file", jsonBodyFile, ts.URL + "/upload"},
			err:    nil,
			output: fmt.Sprintf("JSON request received: %d bytes\n", len(jsonBody)),
		},
		{
			args: []string{"-disable-redirect", ts.URL + "/redirect"},
			err:  errors.New(`Get "/new-url": stopped after 1 redirect`),
		},
		{
			args:   []string{"-header", "Debug-Key1=value1", "-header", "Debug-Key2=value2", ts.URL + "/debug-header-response"},
			err:    nil,
			output: "Debug-Key1=value1 Debug-Key2=value2\n",
		},
		{
			args:   []string{"-basicauth", "user=password", ts.URL + "/debug-basicauth"},
			err:    nil,
			output: "user=password\n",
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

		if len(tc.output) != 0 {
			gotOutput := byteBuf.String()
			if tc.output != gotOutput {
				t.Fatalf("Expected output to be: %s, Got: %s", []byte(tc.output), []byte(gotOutput))
			}
		}
		byteBuf.Reset()
	}
}
