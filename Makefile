.PHONY: test
test:
	go test -v -race -timeout 30s -cover ./...

# generates mocks
.PHONY: generate-mocks
generate-mocks:
	find . -name '*_minimock.go' -delete
	go generate ./...
	go mod tidy -compat=1.17
