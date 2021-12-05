package cmd

import (
	"bytes"
	"errors"
	"net"
	"strings"
	"testing"

	svc "github.com/practicalgo/book-exercise-solutions/chap8/exercise3/service"
	"google.golang.org/grpc"
)

func TestHandleGrpc(t *testing.T) {

	// Setup a real gRPC server
	// we specify the port as 0 so that a port is automatically picked up
	// and then we use the listener's address as the gRPC server address
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	s := grpc.NewServer()
	defer s.Stop()

	svc.RegisterUsersServer(s, &dummyUserService{})
	svc.RegisterRepoServer(s, &dummyReposService{})

	go func() {
		s.Serve(l)
	}()

	testConfigs := []struct {
		args     []string
		err      error
		output   string
		respJson string
	}{
		{
			args: []string{"-service", "Gopher", "-method", "Hello", "-request", "{}", l.Addr().String()},
			err:  errors.New("unrecognized service"),
		},
		{
			args: []string{"-service", "Users", "-method", "GetUser1", "-request", `{"email":"john@doe.com","id":"user-123"}`, l.Addr().String()},
			err:  errors.New("Invalid gRPC method"),
		},
		{
			args:     []string{"-service", "Users", "-method", "GetUser", "-request", `{"email":"john@doe.com","id":"user-123"}`, l.Addr().String()},
			err:      nil,
			respJson: `{"user":{"id":"user-123","firstName":"john","lastName":"doe.com","age":36}}`,
		},
		{
			args: []string{"-service", "Repo", "-method", "GetFoo", "-request", `{"email":"john@doe.com","id":"user-123"}`, l.Addr().String()},
			err:  errors.New("Invalid gRPC method"),
		},
		{
			args:     []string{"-service", "Repo", "-method", "GetRepos", "-request", `{"id":"1"}`, l.Addr().String()},
			err:      nil,
			respJson: `{"repo":[{"id":"repo-123","name":"practicalgo/book-exercise-solutions","url":"git.example.com/practicalgo/book-exercise-solutions","owner":{"id":"user-123"}}]}`,
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range testConfigs {
		err := HandleGrpc(byteBuf, tc.args)
		if tc.err == nil && err != nil {
			t.Fatalf("Expected nil error, got %v", err)
		}

		if tc.err != nil && err == nil {
			t.Fatal("Expected non-nil error, got nil")
		}

		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("Expected error %v, got %v", tc.err, err)
		}

		if len(tc.output) != 0 {
			gotOutput := byteBuf.String()
			if tc.output != gotOutput {
				t.Fatalf("Expected output to be: %#v, Got: %#v", tc.output, gotOutput)
			}
		}

		if len(tc.respJson) != 0 {
			respJson := byteBuf.String()
			// Since the output of protojson.Format() is not stable (as of now), we have to be wary
			// of this breaking
			// As of this writing, i have found that the order of the keys remain stable, but
			// there may be a space between the fields, and sometimes may be not. so before
			// checking the equality, i remove the spaces, sometimes there is also a newline appended
			// hence the expected JSON has been stripped off any newlines or spaces between
			// the fields
			sanitizedRespJson := strings.Replace(string(respJson), " ", "", -1)
			sanitizedRespJson = strings.Replace(string(sanitizedRespJson), "\n", "", -1)

			if sanitizedRespJson != tc.respJson {
				t.Fatalf("Expected result: %v Got: %v", tc.respJson, sanitizedRespJson)
			}
		}

		byteBuf.Reset()
	}
}
