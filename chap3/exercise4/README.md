# Solution to Exercise 3.3

This is my workflow in creating the solution:

- Copy all the code from solution of chapter 3, exercise 2.

In the `cmd` package:

1. Add the following test configurations in `handle_http_test.go`:
```
{
		args:   []string{"-verb", "POST", "-body", "", ts.URL + "/upload"},
		err:    ErrInvalidHTTPPostRequest,
		output: "HTTP POST request must specify a non-empty JSON body\n",
},
{
		args:   []string{"-verb", "POST", "-body", jsonBody, ts.URL + "/upload"},
		err:    nil,
		output: fmt.Sprintf("JSON request received: %d bytes\n", len(jsonBody)),
},
{
		args:   []string{"-verb", "POST", "-body-file", jsonBodyFile, ts.URL + "/upload"},
		err:    nil,
		output: fmt.Sprintf("JSON request received: %d bytes\n", len(jsonBody)),
},
```

Now, the test function will fail for these configurations.

2. Add two new options, `body`  and `body-file` to `httpCmd.go` to have string values
3. If any of this option is specified and a HTTP POST request is made, the data will be sent to the URL as the body
   of a test request
4. For the test server, implement a handler function for `/upload`, as follows:

```
func startTestHttpServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/download", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "this is a response")
	})
	mux.HandleFunc("/upload", func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		data, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "JSON request received: %d bytes", len(data))
	})
	return httptest.NewServer(mux)
}
```

## Trying the final application

With a real server, for example, https://httpbin.org, you should be able to send POST data:

```
C:\> .\exercise3.exe http -verb POST -body  '{"test":"data"}' https://httpbin.org/anything

{
  "args": {},
  "data": "{test:data}",
  "files": {},
  "form": {},
  "headers": {
    "Accept-Encoding": "gzip",
    "Content-Length": "11",
    "Content-Type": "application/json",
    "Host": "httpbin.org",
    "User-Agent": "Go-http-client/2.0",
    "X-Amzn-Trace-Id": "Root=1-6187706b-1f0c77f36d7d16e4447e1dd0"
  },
  "json": null,
  "method": "POST",
  "origin": "61.68.118.26",
  "url": "https://httpbin.org/anything"
}
```