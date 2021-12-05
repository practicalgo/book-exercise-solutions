# Workflow for creating Solution 11.2

- Copy all the files from the code listing, `chap11/pkg-server-2`
- Rename the `TestPackageGetHandler` test function to `TestPackageDownloadHandler`
  - Change the URL path the request is being sent to as `/packages/download?owner_id=1&name=pkg&version=0.1` 
  - Implement a new handler function for handling requets to `/packages/download`
  - Once the test passes, continue to implement the querying functionality
- Rename the package get handler function to package query  handler function
- Update packageQueryResponse type as follows:

```

type pkgQueryResponse struct {
	Packages []pkgRow `json:"packages"`
}
```

- The package query handler function will now return a marshalled version of the above type
as a response
- Add/update tests in package_query_handler_test.go

