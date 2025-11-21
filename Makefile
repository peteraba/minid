default: install

.PHONY: build
build:
	go build -C cmd -o ./build/minid

.PHONY: install
install:
	go build -C cmd -o ~/go/bin/minid

.PHONY: test
test:
	go test ./...

.PHONY: bench
bench:
	go test -bench=. -benchmem ./...