# Solution to Exercise 8.1

This is my workflow in creating the solution:

- Copy all the code from solution of chapter 4, exercise 4.
- Rename go module to: `github.com/practicalgo/book-exercise-solutions/chap8/exercise1` 
  by editing the `go.mod` file and update all package imports accordingly.
- Copy `service` directory from `chap8/user-service`
- Update `grpCmd.go` to  make a call to the users service
- Write tests - i have added example tests to demonstrate the various ways of testing the grpc sub-command
  - `cmd/grpc_flag_parsing_test.go` - this only tests the option parsing behavior
  - `cmd/handle_grpc_users_svc_test.go` - this only tests the behavior when the server only has the users service registered (
    similar to Listing 8.6) and we use a bufconn listener
  - `cmd/handle_grpc_cmd_test.go` - this tests the behavior by creating a "real" gRPC server
  
Usage:

Valid request:

```
C:\> .\exercise1.exe grpc -service Users -method GetUser -request '{\"email\":\"john@doe1.com\",\"id\":\"user-123\"}' localhost:50051
...
```

Invalid request:

```
C:\> .\exercise1.exe grpc -service Users -method GetUser -request '{\"email\":\"john@doe.com\",\"id1\":\"user-123\"}' localhost:50051
proto: (line 1:25): unknown field "id1"
Usage: mync [http|grpc] -h

.. Usage message ..
```

# JSON Pretty printing

```
C:\> .\exercise1.exe grpc -service Users -method GetUser -request '{\"email\":\"john@doe.com\",\"id\":\"user-123\"}' -pretty-print localhost:50051
{
  "user": {
    "id": "user-123",
    "firstName": "john",
    "lastName": "doe.com",
    "age": 36
  }
}
```