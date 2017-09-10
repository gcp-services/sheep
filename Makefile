build:
	go build
	go install

test:
	go test ./...

run:
	${GOPATH}/bin/sheep

.PHONY: test