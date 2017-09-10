build:
	go build
	go install

test:
	go test -covermode=atomic -race ./...

run:
	${GOPATH}/bin/sheep

.PHONY: test