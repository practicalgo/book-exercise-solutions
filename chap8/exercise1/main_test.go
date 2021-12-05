package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
)

var binaryName string
var testHTTPServerURL string
var testGrpcServerURL string

func TestMain(m *testing.M) {
	if runtime.GOOS == "windows" {
		binaryName = "mync.exe"
	} else {
		binaryName = "mync"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "build", "-o", binaryName)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = os.Remove(binaryName)
		if err != nil {
			log.Fatalf("Error removing built binary: %v", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "this is a response")
	})
	mux.HandleFunc("/redirect", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/new-url", http.StatusMovedPermanently)
	})
	ts := httptest.NewServer(mux)
	testHTTPServerURL = ts.URL
	defer ts.Close()

	// Setup a real gRPC server
	// we specify the port as 0 so that a port is automatically picked up
	// and then we use the listener's address as the gRPC server address
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	s := grpc.NewServer()
	defer s.Stop()
	go func() {
		s.Serve(l)
	}()

	testGrpcServerURL = l.Addr().String()
	m.Run()
}

func TestSubcommandInvoke(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	curDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	binaryPath := path.Join(curDir, binaryName)
	t.Log(binaryPath)

	tests := []struct {
		args                []string
		input               string
		expectedOutputLines []string
		expectedExitCode    int
	}{
		{
			args:                []string{},
			expectedExitCode:    1,
			expectedOutputLines: []string{},
		},
		{
			args:             []string{"http"},
			expectedExitCode: 1,
			expectedOutputLines: []string{
				"You have to specify the remote server.",
			},
		},
		{
			args:             []string{"grpc"},
			expectedExitCode: 1,
			expectedOutputLines: []string{
				"You have to specify the remote server.",
			},
		},
		{
			args:             []string{"grpc", "-service", "Gopher", "-method", "Hello", "-request", "{}", testGrpcServerURL},
			expectedExitCode: 1,
			expectedOutputLines: []string{
				"unrecognized service",
			},
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range tests {
		t.Logf("Executing:%v %v\n", binaryPath, tc.args)
		cmd := exec.CommandContext(ctx, binaryPath, tc.args...)
		cmd.Stdout = byteBuf
		if len(tc.input) != 0 {
			cmd.Stdin = strings.NewReader(tc.input)
		}
		err := cmd.Run()

		if err != nil && tc.expectedExitCode == 0 {
			t.Fatalf("Expected application to exit without an error. Got: %v", err)
		}

		if cmd.ProcessState.ExitCode() != tc.expectedExitCode {
			t.Log(byteBuf.String())
			t.Fatalf("Expected application to have exit code: %v. Got: %v", tc.expectedExitCode, cmd.ProcessState.ExitCode())

		}

		output := byteBuf.String()
		lines := strings.Split(output, "\n")
		for num := range tc.expectedOutputLines {
			if lines[num] != tc.expectedOutputLines[num] {
				t.Fatalf("Expected output line to be:%v, Got:%v", tc.expectedOutputLines[num], lines[num])
			}
		}
		byteBuf.Reset()
	}
}
