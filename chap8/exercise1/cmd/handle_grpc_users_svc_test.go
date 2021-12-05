package cmd

import (
	"context"
	"errors"
	"log"
	"net"
	"strings"
	"testing"

	users "github.com/practicalgo/book-exercise-solutions/chap8/exercise1/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type dummyUserService struct {
	users.UnimplementedUsersServer
}

func (s *dummyUserService) GetUser(
	ctx context.Context,
	in *users.UserGetRequest,
) (*users.UserGetReply, error) {
	components := strings.Split(in.Email, "@")
	u := users.User{
		Id:        in.Id,
		FirstName: components[0],
		LastName:  components[1],
		Age:       36,
	}
	return &users.UserGetReply{User: &u}, nil
}

func startTestGrpcServer() (*grpc.Server, *bufconn.Listener) {
	l := bufconn.Listen(10)
	s := grpc.NewServer()
	users.RegisterUsersServer(s, &dummyUserService{})
	go func() {
		err := s.Serve(l)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return s, l
}

func TestCallUsersSvc(t *testing.T) {
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
			c:        grpcConfig{method: "GetUser", request: `{"email":"john@doe.com","id":"user-123"}`},
			err:      nil,
			respJson: `{"user":{"id":"user-123","firstName":"john","lastName":"doe.com","age":36}}`,
		},
		{
			c:        grpcConfig{method: "GetUser", request: "foo-bar"},
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
		usersClient := getUserServiceClient(conn)
		respJson, err := callUsersMethod(usersClient, tc.c)
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
