Exercise 1.1 - functions in the `os/exec` package.
Exercise 2.2 - the `-verb` option. If the verb specifies is anything other than these values, the program should exit with a non-zero exit code and ..

Chapter 6, Page: 11:

"
In the book I write that the innermost middleware is executed first when processing
a request. Unfortunately, it isn't accurate. The request flows from the outermost 
middleware to the innermost middleware in it's journey from the user to the application.
On it's way back, it's the reverse journey - which is what I am describing in the book. 
Since we want to add our request ID before we log or handle panic, we add the new middleware
as the first middleware."



