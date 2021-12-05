package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type wrappedServerStream struct {
	grpc.ServerStream
	messageSent int
	messageRcvd int
}

func (s *wrappedServerStream) SendMsg(m interface{}) error {
	s.messageSent += 1
	return s.ServerStream.SendMsg(m)
}

func (s *wrappedServerStream) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	s.messageRcvd += 1
	return err
}

func loggingUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	resp, err := handler(ctx, req)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Print("No metadata")
	}
	log.Printf("Method:%s, Error:%v, Request-Id:%s",
		info.FullMethod,
		err,
		md.Get("Request-Id"),
	)
	return resp, err
}

func loggingStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {

	err := handler(srv, stream)
	ctx := stream.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Print("No metadata")
	}
	log.Printf("Method:%s, Error:%v, Request-Id:%s",
		info.FullMethod,
		err,
		md.Get("Request-Id"),
	)
	return err
}

func metricUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	end := time.Now()
	log.Printf("Method:%s, Duration:%s",
		info.FullMethod,
		end.Sub(start),
	)
	return resp, err
}

func metricStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	serverStream := &wrappedServerStream{
		ServerStream: stream,
		messageSent:  0,
		messageRcvd:  0,
	}

	start := time.Now()
	err := handler(srv, serverStream)
	end := time.Now()
	log.Printf("Method:%s, Duration:%s, Messages Received:%d, Messages Sent: %d",
		info.FullMethod,
		end.Sub(start),
		serverStream.messageRcvd,
		serverStream.messageSent,
	)
	return err
}
