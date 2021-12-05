# Solution to Exercise 3.1

This is my workflow in creating the solution:

- Copy all the code from solution of chapter 2, exercise 2.

In the `cmd` package:

1. Update the existing test configuration in `handle_http_test.go`:
```
{
		args:   []string{"http://localhost"},
		err:    nil,
		output: "Executing http command\n",
}
```

to be as follows:

```
{
		args:   []string{"http://localhost"},
		err:    nil,
		output: "this is a response\n",
}
```

Now, the test function will fail saying that the expected response doesn't
match the actual response.

2. Copy the definition of `fetchRemoteResource()` function from Listing 3.1 in the book 
   and add it to the file, `httpCmd.go` in the `cmd` package
3. Update the `HandleHttp()` function to now call the `fetchRemoteResource()` function as follows:

```
data, err := fetchRemoteResource(c.url)
if err != nil {
	return err
}
fmt.Fprintln(w, string(data))
return nil
```
4. Now the test will fail with (something like this):

```
Expected nil error, got Get "http://localhost": dial tcp [::1]:80: connectex: No connection could be made because the target machine actively refused it.
```

5. Add a new function to the file, `handle_http_test.go` to create a test HTTP server:

```
func startTestHttpServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/download", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "this is a response")
	})
	return httptest.NewServer(mux)
}
```

6. Update the test config above to be as follows:

```
{
		args:   []string{ts.URL + "/download"},
		err:    nil,
		output: "this is a response\n",
}
```

7. Run the test again, it should pass

8. Now build the application and run it as follows to verify the functionality you implemented:

```
$ go build -o mync  # Linux/MacOS
$ go build -o mync.exe # Windows


C:\> .\mync.exe http -verb GET https://www.github.com
<you should see the response here>
```

In the `main` package:

1. Update `TestMain()` to create a HTTP server:

```
mux := http.NewServeMux()
mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "this is a response")
})
ts := httptest.NewServer(mux)
testServerURL = ts.URL
defer ts.Close()
```

(`testServerURL` is a global variable, declared as `var testServerURL string`)

Update the test configurations to now use testServerURL as the URL. For example:

```
{
		args:             []string{"http", "-verb", "POST", testServerURL},
		expectedExitCode: 0,
		expectedOutputLines: []string{
			"this is a response",
		},
},
```