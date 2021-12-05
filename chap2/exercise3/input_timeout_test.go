package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestInputNoTimeout(t *testing.T) {
	input := strings.NewReader("jane")
	byteBuf := new(bytes.Buffer)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	name, err := getNameContext(ctx, input, byteBuf)

	if err != nil {
		t.Fatalf("Expected nil error, got: %v", err)
	}

	if name != "jane" {
		t.Fatalf("Expected name returned to be jane, got %s", name)
	}
}

func TestInputTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	r, _ := io.Pipe()
	name, err := getNameContext(ctx, r, os.Stdout)
	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	if err == nil {
		t.Fatal("Expected non-nil error, got nil")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Expected error: context.DeadlineExceeded, Got: %s", err)
	}

	if name != "Default Name" {
		t.Fatalf("Expected name returned to be Default Name, got %s", name)
	}
}
