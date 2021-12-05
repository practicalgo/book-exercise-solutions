# Solution to Exercise 5.3

This is my workflow in creating the solution:

1. Initialize a new go module `github.com/practicalgo/book-exercise-solutions/chap5/exercise3`
2. Define the server in server.go

Test the functionality:

```
$ go build
```

Make a request to, `http://localhost:8080/download?filename=go.mod`

**Note**

It's worth noting here, this is in-general vulnerable to path traversal
attacks. Hence, we place a couple of guards around it:

1. We strip away any path separators from the value of `fileName` query parameter. This
restricts the file lookup to happen only in the current directory
2. We also do not allow filenames starting with "."

But, there are likely other ways to exploit this.