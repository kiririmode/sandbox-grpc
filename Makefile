PROTO_DIR=protobuf
DST_DIR=src

prepare:
	brew list --versions protobuf || brew install protobuf
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u google.golang.org/grpc
	go get -u github.com/fullstorydev/grpcurl
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl

compile:
	protoc -I=$(PROTO_DIR) greeter.proto --go_out=plugins=grpc:greeter

.PHONY: prepare compile
