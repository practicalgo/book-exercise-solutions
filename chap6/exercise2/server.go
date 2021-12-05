package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type appConfig struct {
	logger *log.Logger
}

type app struct {
	config  appConfig
	handler func(w http.ResponseWriter, r *http.Request, config appConfig)
}

type requestContextKey struct{}
type requestContextValue struct {
	requestID string
}

func (a app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler(w, r, a.config)
}

func apiHandler(w http.ResponseWriter, r *http.Request, config appConfig) {
	fmt.Fprintf(w, "Hello, world!")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request, config appConfig) {
	if r.Method != "GET" {
		config.logger.Printf("error=\"Invalid request\" path=%s method=%s", r.URL.Path, r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "ok")
}

func panicHandler(w http.ResponseWriter, r *http.Request, config appConfig) {
	panic("I panicked")
}

func setupHandlers(mux *http.ServeMux, config appConfig) {
	mux.Handle("/healthz", &app{config: config, handler: healthCheckHandler})
	mux.Handle("/api", &app{config: config, handler: apiHandler})
	mux.Handle("/panic", &app{config: config, handler: panicHandler})

}

func loggingMiddleware(h http.Handler) http.Handler {
	var requestId string
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)

		ctx := r.Context()
		v := ctx.Value(requestContextKey{})
		if m, ok := v.(requestContextValue); ok {
			requestId = m.requestID
		}
		log.Printf(
			"request_id=%s protocol=%s path=%s method=%s duration=%f",
			requestId, r.Proto, r.URL.Path, r.Method,
			time.Since(startTime).Seconds(),
		)
	})
}

func panicMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rValue := recover(); rValue != nil {
				log.Println("panic detected when handling request", rValue)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Unexpected server error occured")
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func addRequestIdMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.NewString()
		c := requestContextValue{
			requestID: requestID,
		}
		currentCtx := r.Context()
		newCtx := context.WithValue(currentCtx, requestContextKey{}, c)
		rWithContext := r.WithContext(newCtx)
		h.ServeHTTP(w, rWithContext)
	})
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

	m := addRequestIdMiddleware(
		loggingMiddleware(
			panicMiddleware(
				mux,
			),
		),
	)
	log.Fatal(http.ListenAndServe(listenAddr, m))
}
