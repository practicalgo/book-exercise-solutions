# Solution to Exercise 4.3

This is my workflow in creating the solution:

- Copy all the code from solution of chapter 4, exercise 2.
- Rename go module to: `github.com/practicalgo/book-exercise-solutions/chap4/exercise3` 
  by editing the `go.mod` file
- Create a new package `middleware`. Inside it create a new file, `httpLatency.go` with the middleware
  definition. 
- Update HTTP client creation in `cmd\httpCmd.go` to configure the transport

Example:

```
C:\> go build -o mync.exe
C:\> .\mync.exe http -output godev.output https://github.com
2021/11/13 11:25:26 url=https://github.com method=GET protocol=HTTP/1.1 latency=0.085177
Data saved to: godev.output
```