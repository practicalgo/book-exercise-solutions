package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	svc "github.com/practicalgo/code/chap9/bindata-client-streaming/service"
	"google.golang.org/grpc"
)

type appConfig struct {
	filePath  string
	serverUrl string
}

func setupGrpcConn(addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		context.Background(),
		addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
}

func getRepoServiceClient(conn *grpc.ClientConn) svc.RepoClient {
	return svc.NewRepoClient(conn)
}

func uploadRepository(r io.Reader, repoClient svc.RepoClient) (*svc.RepoCreateReply, error) {

	stream, err := repoClient.CreateRepo(
		context.Background(),
	)
	if err != nil {
		return nil, err
	}

	context := svc.RepoCreateRequest_Context{
		Context: &svc.RepoContext{
			CreatorId: "user-123",
			Name:      "test-repo",
		},
	}
	request := svc.RepoCreateRequest{
		Body: &context,
	}
	err = stream.Send(&request)
	if err != nil {
		return nil, err
	}

	// io.Copy uses this so, we just use that too
	size := 32 * 1024
	buf := make([]byte, size)
	for {

		nBytes, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		bData := svc.RepoCreateRequest_Data{
			Data: buf[:nBytes],
		}
		r := svc.RepoCreateRequest{
			Body: &bData,
		}
		err = stream.Send(&r)
		if err != nil {
			return nil, err
		}
	}
	return stream.CloseAndRecv()

}

func setupFlags(w io.Writer, args []string) (appConfig, error) {
	c := appConfig{}
	fs := flag.NewFlagSet("grpc-client", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&c.filePath, "file-path", "", "Repository contents to upload")
	err := fs.Parse(args)
	if err != nil {
		return c, err
	}

	if len(c.filePath) == 0 {
		return c, errors.New("file path empty")
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return c, errors.New("must specify server URL as the only positional argument")
	}
	c.serverUrl = fs.Arg(0)
	return c, nil
}

func main() {
	c, err := setupFlags(os.Stdout, os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	conn, err := setupGrpcConn(c.serverUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	f, err := os.Open(c.filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	repoClient := getRepoServiceClient(conn)
	resp, err := uploadRepository(f, repoClient)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Uploaded %d bytes\n", resp.Size)
}
