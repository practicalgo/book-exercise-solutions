# Solution to Exercise 2.2

This is my workflow in creating the solution:

1. Copy all the code from chap2/sub-cmd-arch
2. Add a new test configuration to `handle_http_test.go`:

```
{
	   args:   []string{"-verb", "PUT", "http://localhost"},
		err:    ErrInvalidHTTPMethod,
		output: "Invalid HTTP method\n",
},
```

3. Add a new function to `httpCmd.go`:

```
func validateConfig(c httpConfig) error {
	allowedVerbs := []string{"GET", "POST", "HEAD"}
	for _, v := range allowedVerbs {
		if c.verb == v {
			return nil
		}
	}
	return ErrInvalidHTTPMethod
}
```

4. In `HandleHTTP()` function, after parsing the flags, call the `validateConfig()` function, check if the error
   returned is due to invalid HTTP verb and display the error message, if so:

```
err = validateConfig(c)
	if err != nil {
		if errors.Is(err, ErrInvalidHTTPMethod) {
			fmt.Fprintln(w, "Invalid HTTP method")
		}
		return err
	}
```

