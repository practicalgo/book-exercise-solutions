module github.com/practicalgo/book-exercise-solutions/chap8/exercise2/client

go 1.17

require google.golang.org/grpc v1.37.0

require (
	github.com/golang/protobuf v1.4.2 // indirect
	golang.org/x/net v0.0.0-20190311183353-d8887717615a // indirect
	golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a // indirect
	golang.org/x/text v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	github.com/practicalgo/book-exercise-solutions/chap8/exercise2/service v0.0.0
	
)

replace github.com/practicalgo/book-exercise-solutions/chap8/exercise2/service => ../service-v1
