# Solution to Exercise 2.1

This is my workflow in creating the solution:

1. Copy all the code from chap2/sub-cmd-arch
2. Create a new file, `main_test.go`
3. Create a `TestMain()` function where I build the application, run the test functions and remove the built application once done
4. In the test function, i then setup test configurations which will invoke the built application with the specified
   sub-command and/or options and verify the exit code as well as the expected output.

Note that when verifying the expected output, I verify the lines of output expected instead of the entire
output.