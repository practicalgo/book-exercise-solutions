module github.com/practicalgo/code/chap9/interceptor-chain/client

go 1.16

require (
	github.com/practicalgo/code/chap9/interceptor-chain/service v0.0.0
	google.golang.org/grpc v1.37.0
)

replace github.com/practicalgo/code/chap9/interceptor-chain/service => ../service
