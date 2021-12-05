package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type appConfig struct {
	logger *log.Logger
}

type app struct {
	config  appConfig
	handler func(w http.ResponseWriter, r *http.Request, config appConfig) (int, error)
}

func (a app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := a.handler(w, r, a.config)
	if err != nil {
		log.Printf("response_status=%d error=%s\n", status, err.Error())
		http.Error(w, err.Error(), status)
		return
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request, config appConfig) (int, error) {
	config.logger.Println("Handling API request")
	fmt.Fprintf(w, "Hello, world!")
	return http.StatusOK, nil
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request, config appConfig) (int, error) {
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, fmt.Errorf("invalid request method:%s", r.Method)
	}
	config.logger.Println("Handling healthcheck request")
	fmt.Fprintf(w, "ok")
	return http.StatusOK, nil
}

func setupHandlers(mux *http.ServeMux, config appConfig) {
	mux.Handle("/healthz", &app{config: config, handler: healthCheckHandler})
	mux.Handle("/api", &app{config: config, handler: apiHandler})
}

func main() {

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	config := appConfig{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
	}

	mux := http.NewServeMux()
	setupHandlers(mux, config)

	log.Fatal(http.ListenAndServe(listenAddr, mux))
}
