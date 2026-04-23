.PHONY: test test-cover test-cover-html

test:
	go test ./...

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func coverage.out

test-cover-html:
	go test ./... -coverprofile=coverage.out
	go tool cover -html coverage.out -o coverage.html
