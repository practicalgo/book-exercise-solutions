package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	users "github.com/practicalgo/book-exercise-solutions/chap8/exercise1/service"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

type grpcConfig struct {
	server      string
	method      string
	request     string
	service     string
	prettyPrint bool
}

func setupGrpcConn(addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		context.Background(),
		addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
}

func getUserServiceClient(conn *grpc.ClientConn) users.UsersClient {
	return users.NewUsersClient(conn)
}

func getUser(
	client users.UsersClient,
	u *users.UserGetRequest,
) (*users.UserGetReply, error) {
	return client.GetUser(context.Background(), u)
}

func createUserRequest(jsonQuery string) (*users.UserGetRequest, error) {
	u := users.UserGetRequest{}
	input := []byte(jsonQuery)
	return &u, protojson.Unmarshal(input, &u)
}

func getUserResponseJson(c grpcConfig, result *users.UserGetReply) ([]byte, error) {
	if c.prettyPrint {
		return []byte(protojson.Format(result)), nil
	}
	return protojson.Marshal(result)

}

func callUsersMethod(usersClient users.UsersClient, c grpcConfig) ([]byte, error) {
	switch c.method {
	case "GetUser":
		req, err := createUserRequest(c.request)
		if err != nil {
			return nil, InvalidInputError{Err: err}
		}
		result, err := getUser(usersClient, req)
		if err != nil {
			return nil, err
		}
		respData, err := getUserResponseJson(c, result)
		return respData, err
	case "":
		return nil, ErrInvalidGrpcMethod
	default:
		return nil, ErrInvalidGrpcMethod
	}
	return nil, nil
}

func validateGrpcConfig(c grpcConfig) error {
	if len(c.service) == 0 {
		return errors.New("unrecognized service")
	}
	if len(c.method) == 0 {
		return ErrInvalidGrpcMethod
	}

	return nil
}

func HandleGrpc(w io.Writer, args []string) error {
	c := grpcConfig{}
	fs := flag.NewFlagSet("grpc", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&c.method, "method", "", "Method to call")
	fs.StringVar(&c.request, "request", "", "Request to send")
	fs.StringVar(&c.service, "service", "", "gRPC service to send the request to")
	fs.BoolVar(&c.prettyPrint, "pretty-print", false, "Pretty print the JSON output")
	fs.Usage = func() {
		var usageString = `
grpc: A gRPC client.

grpc: <options> server`
		fmt.Fprintf(w, usageString)
		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	if err != nil {
		return FlagParsingError{err}
	}
	if fs.NArg() != 1 {
		return InvalidInputError{ErrNoServerSpecified}
	}
	c.server = fs.Arg(0)

	err = validateGrpcConfig(c)
	if err != nil {
		return err
	}

	conn, err := setupGrpcConn(c.server)
	if err != nil {
		return err
	}
	var usersClient users.UsersClient

	switch c.service {
	case "Users":
		usersClient = getUserServiceClient(conn)
		respJson, err := callUsersMethod(usersClient, c)
		if err != nil {
			return err
		}
		fmt.Fprintln(w, string(respJson))
	default:
		return errors.New("unrecognized service")
	}
	return nil
}
