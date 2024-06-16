run:
	docker compose up --build

build:
	docker build -t $appName:runner --target runner .
	go build -o $appName cmd/http/main.go

test:
	docker build -t $appName:tester --target tester .
	docker run $appName:tester

.PHONY: run build test
