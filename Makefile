build:
	go build -o letsgo ./cmd/cli/main/main.go
	chmod +x letsgo

test:
	go test -v ./internal/... 
