package cmd

import (
	"context"
	"errors"
	"net"
	"strings"
	"testing"

	"google.golang.org/grpc"
)

func TestCallReposSvc(t *testing.T) {
	testConfigs := []struct {
		c        grpcConfig
		respJson string
		err      error
	}{
		{
			c:        grpcConfig{},
			err:      ErrInvalidGrpcMethod,
			respJson: "",
		},
		{
			c:        grpcConfig{method: "GetRepos", request: `{"id":"1"}`},
			err:      nil,
			respJson: `{"repo":[{"id":"repo-123","name":"practicalgo/book-exercise-solutions","url":"git.example.com/practicalgo/book-exercise-solutions","owner":{"id":"user-123"}}]}`,
		},
		{
			c:        grpcConfig{method: "GetRepos", request: "foo-bar"},
			err:      InvalidInputError{},
			respJson: "",
		},
	}

	s, l := startTestGrpcServer()
	defer s.GracefulStop()

	bufconnDialer := func(
		ctx context.Context, addr string,
	) (net.Conn, error) {
		return l.Dial()
	}

	conn, err := grpc.DialContext(
		context.Background(),
		"", grpc.WithInsecure(),
		grpc.WithContextDialer(bufconnDialer),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testConfigs {
		t.Log(tc)
		repoClient := getRepoServiceClient(conn)
		respJson, err := callRepoMethod(repoClient, tc.c)
		if tc.err == nil && err != nil {
			t.Fatalf("Expected nil error, got %v", err)
		}

		if tc.err != nil && err == nil {
			t.Fatalf("Expected non-nil error, got nil")
		}

		if tc.err != nil && !errors.As(err, &tc.err) {
			t.Fatalf("Expected error: %v, got: %v", tc.err, err)
		}

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
}
