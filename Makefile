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

# Make a test project in the current directory
project: clean build
	@./build/bin/letsgo make project test-proj repo/module

# Run unit tests on the CLI and the test application
test: build
	@echo "TODO for V2"

# Clean up the build artifacts
clean:
	@rm -rf ./build
	@rm -rf ./test-proj

# Install the binary to the gopath
install: build
	@cp ./build/bin/letsgo $(GOPATH)/bin/letsgo
	@echo "letsgo installed. Run letsgo make something-cool git-host.com/youruser/something-cool"

.PHONY: build project test clean install
