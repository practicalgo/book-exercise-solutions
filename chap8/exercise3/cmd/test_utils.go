package cmd

import (
	"context"
	"log"
	"strings"

	svc "github.com/practicalgo/book-exercise-solutions/chap8/exercise3/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type dummyUserService struct {
	svc.UnimplementedUsersServer
}

func (s *dummyUserService) GetUser(
	ctx context.Context,
	in *svc.UserGetRequest,
) (*svc.UserGetReply, error) {
	components := strings.Split(in.Email, "@")
	u := svc.User{
		Id:        in.Id,
		FirstName: components[0],
		LastName:  components[1],
		Age:       36,
	}
	return &svc.UserGetReply{User: &u}, nil
}

type dummyReposService struct {
	svc.UnimplementedRepoServer
}

func (s *dummyReposService) GetRepos(
	ctx context.Context,
	in *svc.RepoGetRequest,
) (*svc.RepoGetReply, error) {

	repos := []*svc.Repository{
		{
			Id:    "repo-123",
			Name:  "practicalgo/book-exercise-solutions",
			Url:   "git.example.com/practicalgo/book-exercise-solutions",
			Owner: &svc.User{Id: "user-123"},
		},
	}
	return &svc.RepoGetReply{Repo: repos}, nil
}

func startTestGrpcServer() (*grpc.Server, *bufconn.Listener) {
	l := bufconn.Listen(10)
	s := grpc.NewServer()
	svc.RegisterUsersServer(s, &dummyUserService{})
	svc.RegisterRepoServer(s, &dummyReposService{})
	go func() {
		err := s.Serve(l)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return s, l
}
