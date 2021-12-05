# Solution to Exercise 4.4

This is my workflow in creating the solution:

- Copy all the code from solution of chapter 4, exercise 3.
- Rename go module to: `github.com/practicalgo/book-exercise-solutions/chap4/exercise4` 
  by editing the `go.mod` file and update all package imports accordingly.
- Add the new options
- Update httpCmd.go to add the logic for those options
- Updated the validation to exit with an error if output-file is specified with number of requests greater than 1
- The latency middleware is now defined as:

```
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

```

In `httpCmd.go`, the client is now created as:

```
// This is created with all the parameters specified
// when creating http.DefaultTransport, but we configure the 
// MaxIdleConns as per user input
t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          c.maxIdleConns,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	httpLatencyMiddleware := middleware.HttpLatencyClient{
		Logger:    log.New(os.Stdout, "", log.LstdFlags),
		Transport: t,
	}
	httpClient = http.Client{
		CheckRedirect: redirectPolicyFunc,
		Transport:     httpLatencyMiddleware,
	}
```