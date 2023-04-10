# Solution to Exercise 3.4

This is my workflow in creating the solution:

- Copy all the code from solution of chapter 3, exercise 3.

In the `cmd` package:

1. Update the `usageMessage` in `TestHandleHttp` to include two new options,
   `upload` and `form-data`
2. Run the TestHandleHttp test
3. To fix the test, add two new options, `body-file` and `form-data` to accept 
   string values

4. Add the following test configurations to `testConfigs` in `handle_http_test.go`:
```
{
			args: []string{
				"-verb", "POST",
				"-upload", uploadFile,
				"-form-data", "filename=test.data",
				"-form-data", "version=0.1",
				ts.URL + "/upload",
			},
			err: nil,
			output: fmt.Sprintf(
				"HTTP POST request received:filename=test.data,version=0.1,upload=%d bytes",
				len(uploadFile),
			),
		},
		{
			args: []string{
				"-verb", "POST",
				"-body-file", jsonBody,
				"-upload", uploadFile,
				"-form-data", "filename=test.data",
				"-form-data", "version=0.1",
				ts.URL + "/upload",
			},
			err: nil,
			output: fmt.Sprintf(
				"HTTP POST request received:json=%d bytes,filename=test.data,version=0.1,upload=%d bytes",
				len(jsonBody), len(uploadFile),
			),
		},
		{
			args: []string{
				"-verb", "POST",
				"-body", jsonBody,
				"-upload", uploadFile,
				"-form-data", "filename=test.data",
				"-form-data", "version=0.1",
				ts.URL + "/upload",
			},
			err: nil,
			output: fmt.Sprintf(
				"HTTP POST request received:json=%d bytes,filename=test.data,version=0.1,upload=%d bytes",
				len(jsonBody), len(uploadFile),
			),
		},
```

3. We create a file in the temporary directory to upload:

```
uploadData := "This is some data"
	uploadFile := filepath.Join(t.TempDir(), "file.data")
	err = os.WriteFile(uploadFile, []byte(uploadData), 0666)
	if err != nil {
		t.Fatal(err)
	}
```

Now, the test function will fail for these configurations.

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
