.PHONY: clean-test
clean-test:
	go clean -testcache

.PHONY: test
test:
	go test ./... -v -timeout 10s

.PHONY: run
run:
	go run .