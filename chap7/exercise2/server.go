package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func handlePing(w http.ResponseWriter, r *http.Request) {
	log.Println("ping: Got a request")
	time.Sleep(3 * time.Second)
	fmt.Fprintf(w, "pong")
}

func doSomeWork(data []byte) {
	time.Sleep(5 * time.Second)
}

func handleUserAPI(logger *log.Logger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		pingServer := r.URL.Query().Get("ping_server")
		if len(pingServer) == 0 {
			pingServer = "http://localhost:8080"
		}

		done := make(chan bool)

		logger.Println("I started processing the request")

		req, err := http.NewRequestWithContext(
			r.Context(),
			"GET",
			pingServer+"/ping",
			nil,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		client := &http.Client{}
		logger.Println("Outgoing HTTP request")
		resp, err := client.Do(req)
		if err != nil {
			logger.Printf("Error making request: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)

		logger.Println("Processing the response i got")

		go func() {
			doSomeWork(data)
			done <- true
		}()

		select {
		case <-done:
			logger.Println("doSomeWork done: Continuing request processing")
		case <-r.Context().Done():
			logger.Printf("Aborting request processing: %v\n", r.Context().Err())
			return
		}

		fmt.Fprint(w, string(data))
		logger.Println("I finished processing the request")
	}
}

func setupHandlers(mux *http.ServeMux, timeoutDuration time.Duration, logger *log.Logger) {
	userHandler := handleUserAPI(logger)
	hTimeout := http.TimeoutHandler(
		userHandler,
		timeoutDuration,
		"I ran out of time",
	)
	mux.Handle("/api/users/", hTimeout)
	mux.HandleFunc("/ping", handlePing)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	timeoutDuration := 30 * time.Second
	mux := http.NewServeMux()
	setupHandlers(mux, timeoutDuration, logger)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
