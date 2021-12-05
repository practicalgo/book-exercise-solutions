# Solution to Exercise 7.1

This is my workflow in creating the solution:

1. Create a new go module `github.com/practicalgo/book-exercise-solutions/chap7/exercise1`
2. Copy server.go from `chap7/network-request-timeout`
3. Define a new function, `createHTTPGetRequestWithTrace()` to return a `*http.Request` object
   which will make a HTTP GET request and integrates `http.ClientTrace`.
4. Update `handleUserAPI()` function to use the above function to create a `*http.Request` request so
   as to make a HTTP GET request to the `/ping` path.


Build and run the server:

```
$ go build
$ ./exercise1.exe
```

When we make a request to `/api/users/`, we will get a HTTP 503 in the client, and the following
logs in the server:

```
2021/11/22 07:53:13 I started processing the request
2021/11/22 07:53:15 Outgoing HTTP request
2021/11/22 07:53:15 Error making request: Get "http://localhost:8080/ping": context deadline exceeded
```

There is nothing additional as a result of adding the client trace context - which tells us that none
of the HTTP connection setup steps - DNS querying, for example, started.

# Experimentation

Now, update timeoutDuration in `main()` to be 5*time.Second. Build and run the server and make
a request to the `/api/users/` endpoint. You will see the following server logs:

```
2021/11/22 08:00:02 I started processing the request
2021/11/22 08:00:04 Outgoing HTTP request
DNS Start Info: {Host:localhost}
DNS Done Info: {Addrs:[{IP:::1 Zone:} {IP:127.0.0.1 Zone:}] Err:<nil> Coalesced:false}
Got Conn: {Conn:0xc000188008 Reused:false WasIdle:false IdleTime:0s}
2021/11/22 08:00:04 ping: Got a request
Put Idle Conn Error: <nil>
2021/11/22 08:00:04 I finished processing the request
```

1. We first have the DNS record querying and response
2. Then, we get a connection to send the request
3. Then, the ping handler gets the request
4. the ping handler returns a response
5. The connection is put back into the pool


