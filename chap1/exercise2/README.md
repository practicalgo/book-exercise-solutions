# Solution to Exercise 1.2

This is my workflow in creating the solution:

1. Copy all the code from chap1/flag-parse
2. Add a new option, `-o` to the `parseArgs()` function
3. Update validation code to accept num times as 0 when an output path is specified
4. Update tests for `parseArgs()`, `validateArgs()` and `runCmd()` functions
