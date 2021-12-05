# Workflow for creating Solution 11.1

- Copy all the files from the code listing, `chap11/pkg-server-1`
- Add a new test function to package_get_handler_test.go to pass a `?download=true` query parameter
- Update the package get handler function to look for `download` query parameter, and then send the data directly if so