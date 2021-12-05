package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"testing"
	"time"
)

var binaryName string
var testServerURL string

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
	ts := httptest.NewServer(mux)
	testServerURL = ts.URL
	defer ts.Close()
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
			args:             []string{"http", testServerURL},
			expectedExitCode: 0,
			expectedOutputLines: []string{
				"this is a response",
			},
		},
		{
			args:             []string{"http", "-verb", "POST", testServerURL},
			expectedExitCode: 0,
			expectedOutputLines: []string{
				"this is a response",
			},
		},
		{
			args:             []string{"http", "-method", "POST", testServerURL},
			expectedExitCode: 1,
			expectedOutputLines: []string{
				"flag provided but not defined: -method",
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
			args:             []string{"grpc", "127.0.0.1"},
			expectedExitCode: 0,
			expectedOutputLines: []string{
				"Executing grpc command",
			},
		},
		{
			args:             []string{"grpc", "127.0.0.1"},
			expectedExitCode: 0,
			expectedOutputLines: []string{
				"Executing grpc command",
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
