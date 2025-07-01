run: ## run with live reload
	go run ./cmd/server

test: ## run unit tests
	go test -race ./...

vet:
	go vet ./...

lint:
	golangci-lint run

ci: vet test 