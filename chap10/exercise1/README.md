# Workflow for solution to Exercise 10.1

- Copy the `server` and `service` directory from the code listing of Chapter 10, `chap10/server-healthcheck`
- Create a copy of the `server` directory, `server-tls` to initialize a TLS enabled server
- Copy the `tls` directory from the code listing of Chapter 10, `chap10/user-service-tls`

- Create a new directory, `client` and initialize a new module inside it
- Create the command line client, borrowing/copying code from the test for the server healthcheck
- Update the `go.mod` of the client  to contain: 
`replace github.com/practicalgo/code/chap10/server-healthcheck/service => ../service`

- Specify the TLS certificate to the client using the `TLS_CERT_FILE_PATH` environment variable. 
  If that is specified, the client will attempt to create a TLS encrypted connection

  ## Behavior of the client

  - If the healthcheck is successful for the `Check` method, it will print the status, else it will print the error
  - For the `Watch` method, the client will continue running till it gets a non-successful status or an error