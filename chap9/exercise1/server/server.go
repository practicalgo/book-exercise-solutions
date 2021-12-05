package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	svc "github.com/practicalgo/book-exercise-solutions/chap9/exercise1/service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userService struct {
	svc.UnimplementedUsersServer
}

type repoService struct {
	svc.UnimplementedRepoServer
}

func (s *userService) GetUser(
	ctx context.Context,
	in *svc.UserGetRequest,
) (*svc.UserGetReply, error) {
	log.Printf("Received request for user with Email: %s Id: %s\n", in.Email, in.Id)
	components := strings.Split(in.Email, "@")
	if len(components) != 2 {
		return nil, errors.New("invalid email address")
	}
	u := svc.User{
		Id:        in.Id,
		FirstName: components[0],
		LastName:  components[1],
		Age:       36,
	}
	return &svc.UserGetReply{User: &u}, nil
}

func (s *repoService) GetRepos(
	in *svc.RepoGetRequest,
	stream svc.Repo_GetReposServer,
) error {
	log.Printf(
		"Received request for repo with CreateId: %s Id: %s\n",
		in.CreatorId,
		in.Id,
	)
	repo := svc.Repository{
		Owner: &svc.User{Id: in.CreatorId, FirstName: "Jane"},
	}
	cnt := 1
	for {
		repo.Name = fmt.Sprintf("repo-%d", cnt)
		repo.Url = fmt.Sprintf("https://git.example.com/test/%s", repo.Name)
		r := svc.RepoGetReply{
			Repo: &repo,
		}
		if err := stream.Send(&r); err != nil {
			return err
		}
		if cnt >= 5 {
			break
		}
		cnt++
	}
	return nil
}

func (s *repoService) CreateBuild(
	in *svc.Repository,
	stream svc.Repo_CreateBuildServer,
) error {
	log.Printf(
		"Received build request for repo: %s\n", in.Name,
	)
	logLine := svc.RepoBuildLog{
		LogLine:   fmt.Sprintf("Starting build for repository:%s", in.Name),
		Timestamp: timestamppb.Now(),
	}
	if err := stream.Send(&logLine); err != nil {
		return err
	}

	cnt := 1
	for {
		logLine := svc.RepoBuildLog{
			LogLine:   fmt.Sprintf("Build log line  - %d", cnt),
			Timestamp: timestamppb.Now(),
		}
		if err := stream.Send(&logLine); err != nil {
			return err
		}
		if cnt >= 3 {
			break
		}
		time.Sleep(300 * time.Millisecond)
		cnt++
	}
	logLine = svc.RepoBuildLog{
		LogLine:   fmt.Sprintf("Finished build for repository:%s", in.Name),
		Timestamp: timestamppb.Now(),
	}
	if err := stream.Send(&logLine); err != nil {
		return err
	}
	return nil
}

func registerServices(s *grpc.Server) {
	svc.RegisterUsersServer(s, &userService{})
	svc.RegisterRepoServer(s, &repoService{})
}

func startServer(s *grpc.Server, l net.Listener) error {
	return s.Serve(l)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":50051"
	}

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	registerServices(s)
	log.Fatal(startServer(s, lis))
}
