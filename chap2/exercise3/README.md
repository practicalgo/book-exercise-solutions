# Solution to Exercise 2.3

This is my workflow in creating the solution:

1. Copy all the code from chap2/user-input-timeout
2. Refactor the `getNameContext` function to now accept a `io.Reader` and `io.Writer` which will then be
   passed to the `getName` function as arguments
3. Create a new file, `input_timeout_test.go`
4. Create a test function `TestInputNoTimeout` inside it to test the behavior when user input is provided
5. Create a test function, `TestInputTimeout` to test the behavior when user input is not provided

To implement (5), I create a reader, which will not have any data to read, using the [io.Pipe()](https://pkg.go.dev/io#Pipe)
function. I discard the writer, since i don't intend to write anything to it. My initial plan was to simply
use `os.Stdin` instead, but that didn't work since `go test` works by executing the test binary. I was
surprised to find that `go test` wouldn't wait for the user input, so I posted the query
to the golang-nuts google group and my confusion was resolved by another list member.
You can see [the discussion](https://groups.google.com/g/golang-nuts/c/24pL7iQbx64/m/ZHQugkOLAgAJ) on the group.