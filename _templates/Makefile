run:
	go mod tidy
	docker compose up --build

build:
	go mod tidy
	go build -o $appName cmd/http/main.go

test:
	go test -v ./...
