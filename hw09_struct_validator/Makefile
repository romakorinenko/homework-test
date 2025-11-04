.PHONY: test lint

test:
	go test -v -race -count=100 . -cover -v

test_sh:
	sh test.sh


lint:
	golangci-lint run .