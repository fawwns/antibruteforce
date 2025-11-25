.PHONY: run build test clean

run:
	go run ./cmd/antibruteforce/main.go

build:
	go build -o bin/antibruteforce ./cmd/antibruteforce

test:
	go test ./...

clean:
	rm -rf bin/*