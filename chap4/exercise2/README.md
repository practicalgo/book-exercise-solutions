# Solution to Exercise 4.2

This is my workflow in creating the solution:

- Copy all the code from solution of chapter 4, exercise 1.
- Rename go module to: `github.com/practicalgo/book-exercise-solutions/chap4/exercise2` by editing the `go.mod` file


The most challenging aspect for me in this exercise was to be able to add a flag option
which could be specified more than once. This is how we can do that:

```
// we want the user to be able to specify the -header option one or more times
// hence we use this method:
// https://pkg.go.dev/flag#FlagSet.Func
headerOptionFunc := func(v string) error {
	c.headers = append(c.headers, v)
	return nil
}
fs.Func("header", "Add one or more headers to the outgoing request (key=value)", headerOptionFunc)
```

For testing the two new options, I implemented the following two endpoints in the test
http server (handle_http_test.go):

```
mux.HandleFunc("/debug-header-response", func(w http.ResponseWriter, req *http.Request) {
	headers := []string{}
	for k, v := range req.Header {
		if strings.HasPrefix(k, "Debug") {
			headers = append(headers, fmt.Sprintf("%s=%s", k, v[0]))
		}
	}
	fmt.Fprint(w, strings.Join(headers, " "))
})

mux.HandleFunc("/debug-basicauth", func(w http.ResponseWriter, req *http.Request) {
	u, p, ok := req.BasicAuth()
	if !ok {
		http.Error(w, "Basic auth missing/malformed", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "%s=%s", u, p)
})
```