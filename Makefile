build:
	# Build a binary for the CLI
	cp -r _templates ./cmd/cli/
	go build -o letsgo ./cmd/cli/main/main.go
	chmod +x letsgo

build-test:
	# Build a test application to run unit tests against
	make build
	./letsgo make test test

test:
	# Run unit tests on the CLI and the test application
	go test -v ./internal/... 
	make build-test
	(cd test && make test)

clean:
	# Clean up the build artifacts
	rm -rf ./letsgo
	rm -rf ./test

.PHONY: build build-test test clean
