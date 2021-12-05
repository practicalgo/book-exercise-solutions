package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
	healthsvc "google.golang.org/grpc/health/grpc_health_v1"
)

type appConfig struct {
	serverUrl   string
	tlsCertFile string
	service     string
	method      string
}

func setupGrpcConn(c appConfig) (*grpc.ClientConn, error) {

	if len(c.tlsCertFile) == 0 {
		return grpc.DialContext(
			context.Background(),
			c.serverUrl,
			grpc.WithInsecure(),
			grpc.WithBlock(),
		)
	}

	creds, err := credentials.NewClientTLSFromFile(c.tlsCertFile, "")
	if err != nil {
		return nil, err
	}
	credsOption := grpc.WithTransportCredentials(creds)
	return grpc.DialContext(
		context.Background(),
		c.serverUrl,
		credsOption,
		grpc.WithBlock(),
	)
}

func getHealthSvcClient(conn *grpc.ClientConn) healthsvc.HealthClient {
	return healthsvc.NewHealthClient(conn)
}

func callHealthCheckSvc(client healthsvc.HealthClient, c appConfig) error {
	switch c.method {
	case "Check":
		resp, err := client.Check(
			context.Background(),
			&healthsvc.HealthCheckRequest{
				Service: c.service,
			},
		)
		if err != nil {
			return fmt.Errorf("healthcheck call failed:%v", err.Error())
		}
		if resp.Status != healthsvc.HealthCheckResponse_SERVING {
			return fmt.Errorf("%s service not healthy:%v", c.service, resp.Status)
		}
		fmt.Println(resp.Status)
	case "Watch":
		client, err := client.Watch(
			context.Background(),
			&healthsvc.HealthCheckRequest{
				Service: c.service,
			},
		)
		if err != nil {
			log.Fatal(err)
		}

		for {
			resp, err := client.Recv()
			if err != nil {
				return fmt.Errorf("healthcheck call failed:%v", err.Error())
			}
			if resp.Status != healthsvc.HealthCheckResponse_SERVING {
				return fmt.Errorf("%s service not healthy:%v", c.service, resp.Status)
			}
			fmt.Println(resp.Status)
		}
	default:
		log.Fatal("Unrecognized method name")
	}
	return nil
}

func parseFlags(w io.Writer, args []string) (appConfig, error) {
	c := appConfig{}
	fs := flag.NewFlagSet("grpc-health", flag.ContinueOnError)
	fs.SetOutput(w)

	fs.StringVar(&c.service, "service", "", "Service to query the health for")
	fs.StringVar(&c.method, "method", "Check", "Health check method to call - Check/Watch")
	err := fs.Parse(args)
	if err != nil {
		return c, err
	}

	if len(c.service) == 0 {
		return c, errors.New("service name must be non-empty")
	}
	if fs.NArg() != 1 {
		return c, errors.New("must specify gRPC server URL as the only positional argument")
	}
	c.serverUrl = fs.Arg(0)
	return c, nil
}

func main() {

	c, err := parseFlags(os.Stdout, os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	c.tlsCertFile = os.Getenv("TLS_CERT_FILE_PATH")
	conn, err := setupGrpcConn(c)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := getHealthSvcClient(conn)
	if err != nil {
		log.Fatal(err)
	}

	err = callHealthCheckSvc(client, c)
	if err != nil {
		log.Fatal(err)
	}
}
