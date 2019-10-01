GOPATH := $(shell go env GOPATH)
build: gen-rpc
	go build ./...; \
	go install ./...

test: gen-rpc
	go vet ./...; \
	go test -test.short -covermode=atomic ./...

test_acc: gen-rpc
	go vet ./...; \
	go test -covermode=atomic ./...

gen-rpc:
	protoc \
	-I api/v1/ \
	-I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.8.5/third_party/googleapis \
	api/v1/v1.proto --go_out=plugins=grpc:. --grpc-gateway_out=logtostderr=true:. --descriptor_set_out=./api/v1/api_descriptor.pb

run:
	${GOPATH}/bin/sheep

.PHONY: test
