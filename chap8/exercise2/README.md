# Solution to Exercise 8.2

This is my workflow in creating the solution:

## Protocol buffer definitions

- Copy the `service` directory from the book's code repository, `chap8/user-service` into a new directory `service-v1`
- Update the `go.mod` file in `service-v1` directory so that the contents are as follows:

```
module github.com/practicalgo/book-exercise-solutions/chap8/exercise2/service
go 1.17
```
- Copy the `service` directory from the book's code repository, `chap8/user-service` into a new directory `service-v2`
- Update the `go.mod` file in `service-v2` directory so that the contents are as follows:
```
module github.com/practicalgo/book-exercise-solutions/chap8/exercise2/service
go 1.17
```

- Update `users.proto` to add a new field, `location` to the `UserGetReply` message as follows:

```
message UserGetReply {
  User user = 1;
  string location = 2 ;
}
```

- Generate the go code for both `service-v1` and `service-v2`:

```
# Windows
C:\> protoc.exe --go_out=. --go_opt=paths=source_relative `
>> --go-grpc_out=. --go-grpc_opt=paths=source_relative `
>> users.proto

```

## Server using the service-v2 protocol buffer

- Copy the `server` directory from the book's code repository, `chap8/user-service`
- Update the `go.mod` file to use the `service-v2` definitions:
```
module github.com/practicalgo/book-exercise-solutions/chap8/exercise2/server

go 1.17

require google.golang.org/grpc v1.37.1

require github.com/practicalgo/book-exercise-solutions/chap8/exercise2/service v0.0.0

replace github.com/practicalgo/book-exercise-solutions/chap8/exercise2/service => ../service-v2
```
- Update the import statements in the server and the tests to import the service definitions as:
```
users "github.com/practicalgo/book-exercise-solutions/chap8/exercise2/service"
```

- Update the `GetUser()` service handler to now add `Location` to the response:
```
return &users.UserGetReply{User: &u, Location: "Australia"}, nil
```

## Client using the service-v1 protocol buffer

- Copy the `client` directory from the book's code repository, `chap8/user-service`
- Update the `go.mod` file to use the `service-v1` definitions:

```
module github.com/practicalgo/book-exercise-solutions/chap8/exercise2/client

go 1.17

require google.golang.org/grpc v1.37.1
require github.com/practicalgo/book-exercise-solutions/chap8/exercise2/service v0.0.0

replace github.com/practicalgo/book-exercise-solutions/chap8/exercise2/service => ../service-v1
```

- Update the import statements in the client and the tests to import the service definitions as:
```
users "github.com/practicalgo/book-exercise-solutions/chap8/exercise2/service"
```

- Update the `main()` function to print the `User` object:
```
result, err := getUser(
		c,
		&users.UserGetRequest{Email: "jane@doe.com"},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stdout, "%#v\n", result.User)
```

Build, and run the server.

Then, run the client:

```
C:\> .\client.exe localhost:50051
&service.User{state:impl.MessageState{NoUnkeyedLiterals:pragma.NoUnkeyedLiterals{}, DoNotCompare:pragma.DoNotCompare{}, DoNotCopy:pragma.DoNotCopy{}, atomicMessageInfo:(*impl.MessageInfo)(nil)}, sizeCache:0, unknownFields:[]uint8(nil), Id:"", FirstName:"jane", LastName:"doe.com", Age:36}
```

## Client using the service-v2 protocol buffer

- Create a copy of the `client` directory, `client-v2`
- Update the `replace` statement in `go.mod` to point to the `service-v2` directory
- Update `main()` to also print the `Location` field
- Build and run the client:

```
C:\> .\client.exe localhost:50051
&service.User{state:impl.MessageState{NoUnkeyedLiterals:pragma.NoUnkeyedLiterals{}, DoNotCompare:pragma.DoNotCompare{}, DoNotCopy:pragma.DoNotCopy{}, atomicMessageInfo:(*impl.MessageInfo)(nil)}, sizeCache:0, unknownFields:[]uint8(nil), Id:"", FirstName:"jane", LastName:"doe.com", Age:36}: Australia
```
