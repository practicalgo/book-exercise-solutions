# Workflow for Solution 9.1

This is my workflow for creating the solution:

Copy the `service` and `server` directories from the book's code directory, `chap9/server-streaming`

## Updating protocol buffers

- Update the `service/go.mod` file to have the following contents:

```
module github.com/practicalgo/book-exercise-solutions/chap9/exercise1/service

go 1.17
```

- Update `service/repositories.proto` to add the new method and message types:

```
.. other imports ..

import "google/protobuf/timestamp.proto";

service Repo {
  .. previous method ..
  rpc CreateBuild (Repository) returns (stream RepoBuildLog) {}
}

message RepoBuildLog {
  string log_line = 1;
  google.protobuf.Timestamp timestamp = 2;
}

.. rest of the protobuf file ..
```

For timestamps, we store the special `Timestamp` data type as described [here](https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/timestamp.proto). Hence, we  use the import statement at the beginning. The corresponding Go types and guidance
around converting to and from Go types are documented [here](https://pkg.go.dev/google.golang.org/protobuf/types/known/timestamppb).

Generate the update code:

```
C:\> protoc.exe --go_out=. --go_opt=paths=source_relative `
>> --go-grpc_out=. --go-grpc_opt=paths=source_relative `
>> users.proto .\repositories.proto
```

## Update the the repository service

- Update the code in the server and the tests so that they now import the Protobuf generated code using the
import path: `github.com/practicalgo/book-exercise-solutions/chap9/exercise1/service/`
- Update the `go.mod` file to have the following contents:

```
module github.com/practicalgo/book-exercise-solutions/chap9/exercise1/server

go 1.17

require google.golang.org/grpc v1.37.0

require github.com/practicalgo/book-exercise-solutions/chap9/exercise1/service v0.0.0

replace github.com/practicalgo/book-exercise-solutions/chap9/exercise1/service => ../service
```

Ensure your tests pass before moving on

## Adding the new method to the repository service

- Add a new test function to call the `RepoBuild()` method
- Add the method implementation to return 5 log lines, update the test to verify the functionality



