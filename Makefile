dev:
	air
server-run:
	go run ./cmd/server
client-run:
	go run ./cmd/client
fmt:
	golangci-lint fmt
lint:
	golangci-lint run
fix:
	golangci-lint run --fix
clean:
	go clean
