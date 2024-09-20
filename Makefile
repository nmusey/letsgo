.SILENT:

# Build a binary for the CLI
build:
	@mkdir -p ./internal/cli/_templates
	@cp -r _templates ./internal/cli/
	@mkdir -p ./build/bin/
	@go build -o ./build/bin/letsgo ./main.go
	@chmod +x ./build/bin/letsgo
	@echo "letsgo built in ./build/bin"
	@rm -rf ./internal/cli/_templates

# Build a test application to run unit tests against
build-test: clean build
	@./build/bin/letsgo make test-repo test-repo
	@mv ./test-repo ./build/

build-pkg: clean build build-test
	@(cd build/test-repo && ../bin/letsgo pkg test)

# Run a test server for manual testing
run-test: build-test
	@docker compose -f ./build/test-repo/docker-compose.yml up 

# Run unit tests on the CLI and the test application
test: build-test
	@go test ./internal/utils/... # can't test cli package, and don't need to
	@echo -e "\n\033[0;32mAll tests passed\033[m\n"
	@echo -e "\033[0;34mTesting built repo now...\033[m\n"
	@make -C ./build/test-repo test

# Clean up the build artifacts
clean:
	@rm -rf ./build

# Install the binary to the gopath
install: build
	@cp ./build/bin/letsgo $(GOPATH)/bin/letsgo
	@echo "letsgo installed. Run letsgo make something-cool github.com/youruser/something-cool"

.PHONY: build run-test build-test test clean install
