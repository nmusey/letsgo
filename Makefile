build:
	go build -o letsgo ./cmd/cli/main/main.go
	chmod +x letsgo

test:
	go test -v ./internal/... 

clean:
	rm -rf ./letsgo
	rm -rf ./test

test-generated:
	make build
	./letsgo make test test
	(cd ./test && go mod tidy && make test)
	make clean
