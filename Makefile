GOCMD=GO111MODULE=on go

.PHONY: test bench build-all

test:
	$(GOCMD) test -cover -race ./...

bench:
	$(GOCMD) test -bench=. -benchmem ./...

build-all:
	$(GOCMD) build -race ./...
