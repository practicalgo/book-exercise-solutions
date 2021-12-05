# Solution to Exercise 5.1

This is my workflow in creating the solution:

1. Copy all the code from (Listing 5.2 and 5.3) `chap5/http-serve-mux`
2. Update the go module name to: `github.com/practicalgo/book-exercise-solutions/chap5/exercise1`
3. Define a new struct type for storing a log line and a function for emitting the log:

```
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
```

4. Update handler functions to call the `logRequest()` function before processing the request
5. I also added a "catch all" handler function so that requests for path other than `/api` 
   and `/healthcheck` are also logged