build:
	go build -o letsgo ./cmd/cli/main/main.go
	chmod +x letsgo

test:
	go test -v ./internal/... 

test-generated:
	make build
	./letsgo make test test
	(cd ./test && make test)
	rm -rf ./test
