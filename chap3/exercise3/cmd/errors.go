package cmd

import "errors"

var ErrNoServerSpecified = errors.New("You have to specify the remote server.")
var ErrInvalidHTTPMethod = errors.New("Invalid HTTP method")
var ErrInvalidHTTPCommand = errors.New("Invalid HTTP command")
var ErrInvalidHTTPPostCommand = errors.New("Cannot specify both body and body-file")
var ErrInvalidHTTPPostRequest = errors.New("HTTP POST request must specify a non-empty JSON body")
