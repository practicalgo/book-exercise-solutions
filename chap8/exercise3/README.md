# Solution to Exercise 8.3

This is my workflow in creating the solution:

- Copy all the code from solution of chapter 8, exercise 1.
- Rename go module to: `github.com/practicalgo/book-exercise-solutions/chap8/exercise3` 
  by editing the `go.mod` file and update all package imports accordingly.
- Copy `service` directory from the book's code repository `chap8/multiple-services`
- Update `grpCmd.go` to  make a call to the repository service
- Updates tests in:
  - `cmd/grpc_flag_parsing_test.go` - this only tests the option parsing behavior
  - `cmd/handle_grpc_cmd_test.go` - this tests the behavior by creating a "real" gRPC server
- Add new tests to `handle_grpc_repo_svc_test.go`

Usage:

Valid request:

```
C:\> .\exercise3.exe grpc -service Users -method GetUser -request '{\"email\":\"john@doe1.com\",\"id\":\"user-123\"}' localhost:50051
...
```

Invalid request:

```
C:\> .\exercise3.exe grpc -service Users -method GetUser -request '{\"email\":\"john@doe.com\",\"id1\":\"user-123\"}' localhost:50051
proto: (line 1:25): unknown field "id1"
Usage: mync [http|grpc] -h

.. Usage message ..
```

# JSON Pretty printing

```
C:\> .\exercise3.exe grpc -service Users -method GetUser -request '{\"email\":\"john@doe.com\",\"id\":\"user-123\"}' -pretty-print localhost:50051
{
  "user": {
    "id": "user-123",
    "firstName": "john",
    "lastName": "doe.com",
    "age": 36
  }
}
```