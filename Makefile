default: build

.PHONY: build
build:
	go build -C cmd -o ~/go/bin/minid

.PHONY: install
install:
	go mod download

.PHONY: test
test:
	go test ./...