build:
	go build
	go install

web:
	cd web/assets && npm run build

test:
	go vet ./...
	go test -covermode=atomic ./...

run:
	${GOPATH}/bin/sheep

.PHONY: test web