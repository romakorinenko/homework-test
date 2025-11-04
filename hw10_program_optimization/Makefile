.PHONY: test lint

test:
	go test -v -race -count=100 . -cover -v

bench_before:
	go test -v -count=1 -timeout=30s -tags bench . > results/result_before.txt

bench_after:
	go test -v -count=1 -timeout=30s -tags bench . > results/result_after.txt

test_sh:
	sh test.sh

lint:
	golangci-lint run .