# Solution to Exercise 6.1

This is my workflow in creating the solution:

1. Create a new go module `github.com/practicalgo/book-exercise-solutions/chap6/exercise1`
2. Copy server.go from `chap6/http-handler-type`
3. The exercise only mentions that you return an error from the handler functions. However, if you
   only return an error, you end up losing the HTTP status code that should accompany the error response.
   Hence, we update the `app` type as follows:

```
type app struct {
    config  appConfig
	handler func(w http.ResponseWriter, r *http.Request, config appConfig) (int, error)
}
```

In the `ServeHTTP()` method we then process these returned values as follows:

```
status, err := a.handler(w, r, a.config)
if err != nil {
		log.Printf("response_status=%d error=%s\n", status, err.Error())
		http.Error(w, err.Error(), status)
		return
	}
```


4. Demonstration:

Send a POST request to `/healthz` path and you will get back a HTTP 405 error response:

```
$ curl --request POST 192.168.1.109:8080/healthz
invalid request method:POST
```

On the server, you will see the following log lines:

```
2021/11/20 17:39:28 server.go:38: Handling healthcheck request
2021/11/20 17:39:46 response_status=405 error=invalid request method:POST
```