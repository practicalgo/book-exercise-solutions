package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_DecodeHandler(t *testing.T) {
	const jsonStream = `
	{"user_ip": "172.121.19.21", "event": "click_on_add_cart"}{"user_ip": "172.121.19.21", "event": "click_on_checkout"}
`
	body := strings.NewReader(jsonStream)

	r := httptest.NewRequest("POST", "http://example.com/decode", body)
	w := httptest.NewRecorder()

	decodeHandler(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("Expected Response Status: %v, Got: %v", http.StatusOK, w.Result().StatusCode)
	}
}

func Test_DecodeUnknownFieldError(t *testing.T) {
	const jsonStream = `
	{"user_ip": "172.121.19.21", "event": "click_on_add_cart", "user_data":"some_data"}{"user_ip": "172.121.19.21", "event": "click_on_checkout"}
`
	body := strings.NewReader(jsonStream)

	r := httptest.NewRequest("POST", "http://example.com/decode", body)
	w := httptest.NewRecorder()

	decodeHandler(w, r)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected Response Status: %v, Got: %v", http.StatusBadRequest, w.Result().StatusCode)
	}
}
