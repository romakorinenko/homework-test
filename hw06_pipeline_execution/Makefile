.PHONY: test lint

test:
	go test -v -race -count=100 . -cover -v


lint:
	golangci-lint run .