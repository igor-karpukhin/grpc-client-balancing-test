.PHONY: all clean proto client server

all: proto

proto:
	protoc --go_out=plugins=grpc:${GOPATH}/src grpc/*.proto

server:
	go build server/*.go -o server

client:
	go build client/*.go -o client

clean:
	rm server client
