package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func metadataUnaryInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	ctxWithMetadata := metadata.AppendToOutgoingContext(
		ctx,
		"Request-Id", "request-123",
	)
	err := invoker(ctxWithMetadata, method, req, reply, cc, opts...)
	return err
}

func metadataStreamingInterceptor(
	ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	ctxWithMetadata := metadata.AppendToOutgoingContext(
		ctx,
		"Request-Id", "request-123",
	)
	clientStream, err := streamer(ctxWithMetadata, desc, cc, method, opts...)
	return clientStream, err
}

func loggingUnaryInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	end := time.Now()

	log.Printf("Method:%s, Duration:%s, Error:%v",
		method,
		end.Sub(start),
		err,
	)
	return err
}

type wrappedClientStream struct {
	grpc.ClientStream
	messageSent int
	messageRcvd int
}

func (s *wrappedClientStream) SendMsg(m interface{}) error {
	log.Printf("Send msg called: %T", m)
	s.messageSent += 1
	return s.ClientStream.SendMsg(m)
}

func (s *wrappedClientStream) RecvMsg(m interface{}) error {
	log.Printf("Recv msg called: %T", m)
	err := s.ClientStream.RecvMsg(m)
	s.messageRcvd += 1
	return err
}

func (s *wrappedClientStream) CloseSend() error {
	log.Println("CloseSend() called")
	v := s.Context().Value(streamDurationContextKey{})

	if m, ok := v.(streamDurationContextValue); ok {
		log.Printf("Duration:%v", time.Since(m.startTime))
	}
	err := s.ClientStream.CloseSend()
	log.Printf("Messages Sent: %d, Messages Received:%d\n",
		s.messageSent,
		s.messageRcvd,
	)
	return err
}

type streamDurationContextKey struct{}
type streamDurationContextValue struct {
	startTime time.Time
}

func loggingStreamingInterceptor(
	ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	c := streamDurationContextValue{
		startTime: time.Now(),
	}
	ctxWithTimer := context.WithValue(
		ctx,
		streamDurationContextKey{},
		c,
	)

	stream, err := streamer(
		ctxWithTimer,
		desc,
		cc, method,
		opts...,
	)
	clientStream := &wrappedClientStream{
		ClientStream: stream,
		messageRcvd:  0,
		messageSent:  0,
	}
	return clientStream, err
}
