build:
	CGO_ENABLED=0 go build -o scout ./cmd/scout/

run: build
	./scout

test:
	go test ./...

clean:
	rm -f scout

.PHONY: build run test clean
