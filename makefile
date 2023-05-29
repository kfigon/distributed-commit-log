.PHONY: clean-test
clean-test:
	go clean -testcache

.PHONY: test-verbose
test-verbose:
	go test ./... -v -timeout 10s

.PHONY: test
test:
	go test ./... -timeout 10s

.PHONY: coverage
coverage:
	go test ./... -cover -timeout 10s

.PHONY: run
run:
	go run .