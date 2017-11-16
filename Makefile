build:
	${GOPATH}/bin/statik -src=web/assets/dist
	go build
	go install

web:
	cd web/assets && npm run build
	${GOPATH}/bin/statik -src=web/assets/dist

test:
	go vet ./...
	go test -covermode=atomic ./...

run:
	${GOPATH}/bin/sheep

.PHONY: test web