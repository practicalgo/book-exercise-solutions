package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type logLine struct {
	URL           string `json:"url"`
	Method        string `json:"method"`
	ContentLength int64  `json:"content_length"`
	Protocol      string `json:"protocol"`
}

func logRequest(req *http.Request) {
	l := logLine{
		URL:           req.URL.String(),
		Method:        req.Method,
		ContentLength: req.ContentLength,
		Protocol:      req.Proto,
	}
	data, err := json.Marshal(&l)
	if err != nil {
		panic(err)
	}
	log.Println(string(data))
}

func apiHandler(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	fmt.Fprintf(w, "Hello, world!")
}

func healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	fmt.Fprintf(w, "ok")
}

func catchAllHandler(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	fmt.Fprintf(w, "your request was processed by the catch all handler")
}

func setupHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", healthCheckHandler)
	mux.HandleFunc("/api", apiHandler)
	mux.HandleFunc("/", catchAllHandler)
}

func main() {

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	setupHandlers(mux)

	log.Fatal(http.ListenAndServe(listenAddr, mux))
}
