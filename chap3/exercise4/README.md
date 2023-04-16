# Solution to Exercise 3.4

The solution builds on top of the solution to [exercise 3.3](../exercise3/)

## Multiple values for a flag option

To accept multiple values for the `-form-data` flag option, we use the technique
described in [my blog post](https://echorand.me/posts/go-flag-option-append/).
This technique is also used in solutions to exercises in [chap 4](../../chap4/)
and [chap 8](../../chap8/).

## Sending JSON data in a multipart/form-data body

The exercise requires us to find a way to send JSON data if specified
by the program's user along with any files and form data. We adopt a technique
that seems to be commonly used for such a use case:

- If any JSON data is specified using the `body` or `body-file` option (Exercise 3.3)
  along with `-upload` and/or `form-data`, we add a specially named form field,
  `jsondata` to contain the JSON data. The request body thus will continue to have
  the `Content-Type` set to `multipart/form-data`
- On the server side, we know if we get any data in the `jsondata` form field, we
  need to unmarshal the JSON data appropriately

## Using golden files for testing

We also start using "golden files" for testing certain behavior of our solution.
This technique is described in this [blog post](https://ieftimov.com/posts/testing-in-go-golden-files),
however, we don't follow it exactly as described in the post. 

The [testdata](./testdata/) directory contains the following files:

- expectedGolden.0: Expected output for the test case, "0" in the `main` package
- expectedGolden.1: Expected output for the test case, "1" in the `main` package
- expectedGolden.2: Expected output for the test case, "2" in the `main` package
- expectedGolden.cmd.httpCmdUsage:  Expected output for the only test case where we
  use this technique in the `cmd` package

  If we update our code, and a test fails, it will write the expected output to
  a file, "gotOutput.xx.xx" which you are then required to manually inspect
  and overwrite the corresponding golden file in `testdata`.


## Things to improve

- Expected/golden logging is hard to read
- Test configs need a good identification, currently have to manually look at the array indexes
