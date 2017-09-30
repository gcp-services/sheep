build:
	go build
	go install

test:
	go test -covermode=atomic ./...

run:
	${GOPATH}/bin/sheep

.PHONY: test