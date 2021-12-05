# Solution to Exercise 6.2

This is my workflow in creating the solution:

1. Create a new go module `github.com/practicalgo/book-exercise-solutions/chap6/exercise2`
2. Copy server.go from `chap6/middleware-chaining`
3. Update the logging middleware to log the request ID
4. Add a new middleware to add the request ID
5. I have used an external package: `github.com/google/uuid` to create a New UUID as a string and add that
   as the request ID
6. Update the `main()` function to also add the request middleware to the chain.

In the book I write that the innermost middleware is executed first when processing
a request. Unfortunately, it isn't accurate. The request flows from the outermost 
middleware to the innermost middleware in it's journey from the user to the application.
On it's way back, it's the reverse journey - which is what I am describing in the book. 
Since we want to add our request ID before we log or handle panic, we add the new middleware
as the first middleware.