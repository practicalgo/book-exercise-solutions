# Workflow for solution to Exercise 9.3

- Copy all of the contents of the directory from the book's source repository, `chap9/interceptor-chain`

## Updating the server interceptor

- update `wrappedServerStream` to add two new fields, `messageSent` and `messageRecvd`:

```
type wrappedServerStream struct {
	grpc.ServerStream
	messageSent int
	messageRcvd int
}
```

- Update the SendMsg() and RecvMsg() methods to now have pointer receivers and then increment the counters
- Update the logging interceptor to not use the wrapped stream
- Update the metric interceptor to wrap the incoming stream inside a `wrappedStream` value:

```
serverStream := &wrappedServerStream{
		ServerStream: stream,
		messageSent:  0,
		messageRcvd:  0,
}
```

## Updating the client interceptor

- update `wrappedClientStream` to add two new fields, `messageSent` and `messageRecvd`:

```
type wrappedClientStream struct {
	grpc.ClientStream
	messageSent int
	messageRcvd int
}
```

- Update the SendMsg() and RecvMsg() methods to now have pointer receivers and then increment the counters
- Update the `loggingStreamingInterceptor` to initialize the `wrappedClientStream` as follows:

```
clientStream := &wrappedClientStream{
		ClientStream: stream,
		messageRcvd:  0,
		messageSent:  0,
	}
```
- Finally update the CloseSend()  method of the interceptor to log the messages received and sent count:

```
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
```