package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/practicalgo/book-exercise-solutions/chap4/exercise4/cmd"
)

var errInvalidSubCommand = errors.New("Invalid sub-command specified")

func printUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: mync [http|grpc] -h\n")
	cmd.HandleHttp(w, []string{"-h"})
	cmd.HandleGrpc(w, []string{"-h"})
}

func handleCommand(w io.Writer, args []string) error {
	var err error

	if len(args) < 1 {
		err = cmd.InvalidInputError{Err: errInvalidSubCommand}
	} else {
		switch args[0] {
		case "http":
			err = cmd.HandleHttp(w, args[1:])
		case "grpc":
			err = cmd.HandleGrpc(w, args[1:])
		case "-h":
			printUsage(w)
		case "-help":
			printUsage(w)
		default:
			err = cmd.InvalidInputError{Err: errInvalidSubCommand}
		}
	}

	// for non-nil errors, we can have three kinds of errors
	// 1. Flag parsing error (E.g. use of invalid option)
	// 2. Invalid input error (E.g. JSON body specified for a GET request)
	// 3. Application specific error (E.g. remote server returned an error for example)
	// For (1), the flag package will show the error and also print the usage, so we don't do anything here
	// For (2), the we want to show the error and print the usage of the program
	// For (3), we only want to show the error
	if err != nil {
		if !errors.As(err, &cmd.FlagParsingError{}) {
			fmt.Fprintln(w, err.Error())
		}
		if errors.As(err, &cmd.InvalidInputError{}) {
			printUsage(w)
		}
	}
	return err
}

func main() {
	err := handleCommand(os.Stdout, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}
}
