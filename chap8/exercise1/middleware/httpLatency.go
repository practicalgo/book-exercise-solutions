package middleware

import (
	"log"
	"net/http"
	"time"
)

type HttpLatencyClient struct {
	Logger    *log.Logger
	Transport http.RoundTripper
}

func (c HttpLatencyClient) RoundTrip(
	r *http.Request,
) (*http.Response, error) {
	startTime := time.Now()
	resp, err := c.Transport.RoundTrip(r)
	c.Logger.Printf(
		"url=%s method=%s protocol=%s latency=%f\n",
		r.URL, r.Method, r.Proto, time.Since(startTime).Seconds(),
	)
	return resp, err
}
