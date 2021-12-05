# Worklow for Solution 9.2

This is my workflow for creating the solution.

- Copy `server` and `service` directories from the book's code repository, `chap9/bindata-client-streaming`

- Create a new `client` directory where we will create our client
- Initialize a new module,  `go mod init github.com/practicalgo/book-exercise-solutions/chap9/exercise1/client`
- Use the import path for the service definitions as we did in chap9/ code:

```
svc "github.com/practicalgo/code/chap9/bindata-client-streaming/service"
```

The `go.mod` for the client will be as follows:

```
module github.com/practicalgo/book-exercise-solutions/chap9/exercise1/client

require (
	github.com/practicalgo/code/chap9/bindata-client-streaming/service v0.0.0
	google.golang.org/grpc v1.42.0
)


replace github.com/practicalgo/code/chap9/bindata-client-streaming/service => ../service

go 1.17

```


## Demonstration

Run the server and then run the client.

```
 .\client.exe -file-path .\client.exe localhost:50051
Uploaded 11435520 bytes
```

(The exercise mentions a tar.gz file, but of course we can use any file for testing the upload behavior)