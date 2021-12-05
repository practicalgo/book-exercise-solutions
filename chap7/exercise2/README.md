# Solution to Exercise 7.2

This is my workflow in creating the solution:

1. Create a new go module `github.com/practicalgo/book-exercise-solutions/chap7/exercise1`
2. Copy server.go from `chap7/client-disconnect-handling`
3. Refactor the server code so that the users API handler now accepts a configurable logger as an argument
4. Also, updated the user api handler to accept the address of the server for the ping handler as a query parameter
4. Write the test with a client timeout configured and verify that the logs are as expected