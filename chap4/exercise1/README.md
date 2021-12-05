# Solution to Exercise 4.1

This is my workflow in creating the solution:

- Copy all the code from solution of chapter 3, exercise 3.
- Rename go module to: `github.com/practicalgo/book-exercise-solutions/chap4/exercise1` by editing the `go.mod` file
- Add the following test configuration to `cmd/handle_http_test.go`:

```
{
			args: []string{"-disable-redirect", ts.URL + "/redirect"},
			err:  errors.New(`Get "/new-url": stopped after 1 redirect`),
},
```
- The test will fail now since the `-disable-redirect` option is not recognized
- Add a new bool option, disable-redirect to `http` flag in `cmd/httpCmd.go`
- Implement a new handler in the test HTTP server:

```
mux.HandleFunc("/redirect", func(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/new-url", http.StatusMovedPermanently)
})
```

The above handler will return a redirect to http://127.0.0.1:<random-port>/new-url as a response
to a request made to the "/redirect" path.

- Ensure that the test still fails, as the client has not been configured to not follow redirects
- We will now implement support in the HTTP flag to disable the default redirection behavior
- Ensure tests pass
- Next, update the tests in the `main` package and add a new test config to `main_test.go`:

```
{
		args:             []string{"http", "-disable-redirect", testServerURL},
		expectedExitCode: 1,
		expectedOutputLines: []string{
			`Get "/new-url": stopped after 1 redirect`,
		},
},
```

- In Chapter 1, we used the `errors.Is()` function to display an error message in the `main()` 
function when the error is a non flag parsing error. Here we do slighly differently.
We create two new types in the `cmd` package:

```
type FlagParsingError struct {
	err error
}

type InvalidInputError struct {
	Err error
}
```

Then, we update the `httpCmd.go` file to now return errors which are due to bad input as follows
(for example):

```
if len(postBodyFile) != 0 && len(c.postBody) != 0 {
		return InvalidInputError{ErrInvalidHTTPPostCommand}
	}
```

Then, in `handleCommand()` function, we use the `errors.As()` function to check the
error typpe and hande them differently:

```
if err != nil {
		if !errors.As(err, &cmd.FlagParsingError{}) {
			fmt.Fprintln(w, err.Error())
		}
		if errors.As(err, &cmd.InvalidInputError{}) {
			printUsage(w)
		}
	}
```

If the error is due to an invalid input, then we also print the usage message, else, we only
display the error.

After the changes, ensure all the tests pass.  

Then, run the `http` sub-command to verify the functionality of disable-redirect:

```
C:\> go build -o mync.exe

C:\> .\mync.exe http -disable-redirect http://www.github.com
Get "https://www.github.com/": stopped after 1 redirect
```

