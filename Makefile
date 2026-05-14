run:
	go run ./cmd/gitviz

build:
	go build -o bin/gitviz ./cmd/gitviz

test:
	go test ./...

fmt:
	go fmt ./...

lint:
	go vet ./...
