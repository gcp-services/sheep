build:
	go build
	go install

test:
	go vet ./...
	go test -covermode=atomic ./...

run:
	${GOPATH}/bin/sheep

.PHONY: test