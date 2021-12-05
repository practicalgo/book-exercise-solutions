# Solution to Exercise 6.3

This is my workflow in creating the solution:

1. Create a new go module `github.com/practicalgo/book-exercise-solutions/chap6/exercise3`
2. Copy all the `.go` files (including the subdirectories) from `chap6/complex-server`)
3. Add a new test function to `handlers/handlers_test.go` for the health check handler
4. Since we want to test the behavior of the handler function for multiple HTTP request methods,
   we adopt a table driven testing approach and create a list of test configurations as follows:

   ```
   testConfigs := []struct {
		httpMethod           string
		expectedStatus       int
		expectedResponseBody string
	}{
		{
			httpMethod:           "GET",
			expectedStatus:       http.StatusOK,
			expectedResponseBody: "ok",
		},
		{
			httpMethod:           "POST",
			expectedStatus:       http.StatusMethodNotAllowed,
			expectedResponseBody: "Method not allowed\n",
		},
		{
			httpMethod:           "PUT",
			expectedStatus:       http.StatusMethodNotAllowed,
			expectedResponseBody: "Method not allowed\n",
		},
	}
    ```