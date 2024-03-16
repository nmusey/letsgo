run:
	docker compose up --build

build:
	go build -o $appName cmd/http/main.go

test:
	docker build -t $appName:tester --target tester .
	docker run $appName:tester

.PHONY: run build test
