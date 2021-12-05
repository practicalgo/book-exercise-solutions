package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"testing"

	svc "github.com/practicalgo/book-exercise-solutions/chap9/exercise1/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func startTestGrpcServer() *bufconn.Listener {
	l := bufconn.Listen(10)
	s := grpc.NewServer()
	registerServices(s)
	go func() {
		log.Fatal(startServer(s, l))
	}()
	return l
}
func TestUserService(t *testing.T) {

	l := startTestGrpcServer()

	bufconnDialer := func(
		ctx context.Context, addr string,
	) (net.Conn, error) {
		return l.Dial()
	}

	client, err := grpc.DialContext(
		context.Background(),
		"", grpc.WithInsecure(),
		grpc.WithContextDialer(bufconnDialer),
	)
	if err != nil {
		t.Fatal(err)
	}
	usersClient := svc.NewUsersClient(client)
	resp, err := usersClient.GetUser(
		context.Background(),
		&svc.UserGetRequest{Id: "foo-bar", Email: "jane@doe.com"},
	)

	if err != nil {
		t.Fatal(err)
	}
	if resp.User.FirstName != "jane" {
		t.Errorf(
			"Expected FirstName to be: jane, Got: %s",
			resp.User.FirstName,
		)
	}

}

func TestRepoService(t *testing.T) {

	l := startTestGrpcServer()

	bufconnDialer := func(
		ctx context.Context, addr string,
	) (net.Conn, error) {
		return l.Dial()
	}

	client, err := grpc.DialContext(
		context.Background(),
		"", grpc.WithInsecure(),
		grpc.WithContextDialer(bufconnDialer),
	)
	if err != nil {
		t.Fatal(err)
	}
	repoClient := svc.NewRepoClient(client)
	stream, err := repoClient.GetRepos(
		context.Background(),
		&svc.RepoGetRequest{CreatorId: "user-123", Id: "repo-123"},
	)
	if err != nil {
		t.Fatal(err)
	}
	var repos []*svc.Repository
	for {
		repo, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		repos = append(repos, repo.Repo)
	}
	if len(repos) != 5 {
		t.Fatalf("Expected to get back 5 repos, got back: %d repos", len(repos))
	}

	for idx, repo := range repos {
		gotRepoName := repo.Name
		expectedRepoName := fmt.Sprintf("repo-%d", idx+1)

		if gotRepoName != expectedRepoName {
			t.Errorf(
				"Expected Repo Name to be: %s, Got: %s",
				expectedRepoName,
				gotRepoName,
			)
		}
	}
}

func TestRepoBuildMethod(t *testing.T) {

	l := startTestGrpcServer()

	bufconnDialer := func(
		ctx context.Context, addr string,
	) (net.Conn, error) {
		return l.Dial()
	}

	client, err := grpc.DialContext(
		context.Background(),
		"", grpc.WithInsecure(),
		grpc.WithContextDialer(bufconnDialer),
	)
	if err != nil {
		t.Fatal(err)
	}
	repoClient := svc.NewRepoClient(client)
	stream, err := repoClient.CreateBuild(
		context.Background(),
		&svc.Repository{Name: "practicalgo/test-repo"},
	)
	if err != nil {
		t.Fatal(err)
	}
	var logLines []*svc.RepoBuildLog
	for {
		line, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		logLines = append(logLines, line)
	}
	if len(logLines) != 5 {
		t.Fatalf("Expected to get back 3 lines in the log, got back: %d repos", len(logLines))
	}

	// assert the first and last lines
	expectedFirstLine := "Starting build for repository:practicalgo/test-repo"
	if logLines[0].LogLine != expectedFirstLine {
		t.Fatalf("Expected first line to be:%s, Got:%s", expectedFirstLine, logLines[0].LogLine)
	}
	expectedLastLine := "Finished build for repository:practicalgo/test-repo"
	if logLines[4].LogLine != expectedLastLine {
		t.Fatalf("Expected last line to be:%s,Got:%s", expectedLastLine, logLines[4].LogLine)
	}

	// Ensure that we can unmarshal the incoming timestamps into a Go time.Time type for
	// any of the log lines
	// https://pkg.go.dev/google.golang.org/protobuf/types/known/timestamppb
	logLine := logLines[0]
	if err := logLine.Timestamp.CheckValid(); err != nil {
		t.Fatalf("Logline timestamp invalid: %#v", logLine)
	}
}
