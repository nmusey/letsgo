.SILENT:

# Build a binary for the CLI
build:
	@mkdir -p ./build/src/
	@mkdir -p ./build/bin/
	@mkdir -p ./internal/cli/_templates
	@cp -r _templates ./internal/cli/_templates
	@cp main.go ./build/src/main.go
	@go build -o ./build/bin/letsgo ./build/src/main.go
	@chmod +x ./build/bin/letsgo
	@rm -rf ./internal/cli/_templates
	@echo "letsgo built in ./build/bin"

# Build a test application to run unit tests against
build-test: clean build
	@./build/bin/letsgo make test-repo test-repo
	@mv ./test-repo ./build/

run-test: build-test
	# Run a test server for manual testing
	@make build-test
	@docker compose -f ./build/test-repo/docker-compose.yml up 

test: build-test
	# Run unit tests on the CLI and the test application
	@go test ./internal/utils/... # can't test cli package, and don't need to
	@echo -e "\n\033[0;32mAll tests passed\033[m\n"
	@echo -e "\033[0;34mTesting built repo now...\033[m\n"
	@make build-test
	@make -C ./build/test-repo test

clean:
	# Clean up the build artifacts
	@rm -rf ./build

install: build
	@cp ./build/bin/letsgo $(GOPATH)/bin/letsgo
	@echo "letsgo installed. Run letsgo make something-cool github.com/youruser/something-cool"

.PHONY: build run-test build-test test clean install
