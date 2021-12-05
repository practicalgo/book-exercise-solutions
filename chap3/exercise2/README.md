# Solution to Exercise 3.2

This is my workflow in creating the solution:

- Copy all the code from solution of chapter 3, exercise 1.

In the `cmd` package:

1. Add a new test configuration in `handle_http_test.go`:
```
{
		args:   []string{"-verb", "GET", "-output", "file_path.out", "http://localhost/binary-data"},
		err:    nil,
		output: "Data saved to file_path.out",
},
```

Now, the test function will fail.

2. Add a new option, `output` to `httpCmd.go` to have a string value
3. If this option is specified, the response from the server will now be written to a file and
   if that completes successfully, display a message to the user, "Data saved to: <file path>".
