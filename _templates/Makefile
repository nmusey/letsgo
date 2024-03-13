run:
	go mod tidy
	docker compose up --build

build:
	go mod tidy
	go build -o $appName cmd/http/main.go

test:
	docker build -t $appName:tester --target tester .
	docker run $appName:tester

test-debug:
	go test -v ./...

.PHONY: run build test test-debug
