# Solution to Exercise 1.1

This is my workflow in creating the solution:

1. Copy all the code from chap1/manual-parse
2. Create a new file, `main_test.go`
3. Create a `TestMain()` function where I build the application, run the test functions and remove the built application once done
4. In the test function, i then setup test configurations which will invoke the built application with the specified
   command line arguments and/or input and verify the exit code as well as the expected output.

Note that when verifying the expected output, I verify the lines of output expected instead of the entire
output.