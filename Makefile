version := $(shell git tag | tail -n1)
commit := $(shell git rev-parse HEAD)

withenv: main.go
	go build -ldflags "-w -s -X main.version=${version} -X main.commit=${commit}" -o $@
