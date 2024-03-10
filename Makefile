build:
	# Build a binary for the CLI
	cp -r _templates ./cmd/cli/
	go build -o letsgo ./cmd/cli/main/main.go
	chmod +x letsgo

build-test:
	# Build a test application to run unit tests against
	make build
	./letsgo make test test
	(cd test && templ generate)

run-test:
	# Run a test server for manual testing
	make build-test
	(cd test && docker compose up)
	make clean

test:
	# Run unit tests on the CLI and the test application
	go test ./internal/... 
	make build-test
	(cd test && make test)

test-debug:
	go test -v ./internal/...
	make build-test
	(cd test && make test-debug)

clean:
	# Clean up the build artifacts
	rm -rf ./letsgo
	rm -rf ./test
	rm -rf ./cmd/cli/_templates

.PHONY: build build-test test test-debug clean
