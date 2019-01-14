
build: gen-rpc
	go build ./...; \
	go install ./...

test: gen-rpc
	go vet ./...; \
	go test -covermode=atomic ./...

gen-rpc:
	protoc \
	-I api/v1/ \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	api/v1/v1.proto --go_out=plugins=grpc:${GOPATH}/src --grpc-gateway_out=logtostderr=true:${GOPATH}/src
run:
	${GOPATH}/bin/sheep

.PHONY: test
