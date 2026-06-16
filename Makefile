.SILENT:

# Build a binary for the CLI
build:
	@mkdir -p ./build/bin/
	@go build -o ./build/bin/letsgo ./main.go
	@chmod +x ./build/bin/letsgo
	@echo "letsgo built in ./build/bin"

# Creates a release of the built application
release: build
	@echo "TODO for V2: GH action with release"		

# Run a fresh binary
run: build
	@./build/bin/letsgo

# Run unit tests on the CLI and the test application
test: build
	@echo "TODO for V2"

# Clean up the build artifacts
clean:
	@rm -rf ./build

# Install the binary to the gopath
install: build
	@cp ./build/bin/letsgo $(GOPATH)/bin/letsgo
	@echo "letsgo installed. Run letsgo make something-cool github.com/youruser/something-cool"

.PHONY: build run-test build-test test clean install
