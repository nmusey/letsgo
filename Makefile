build:
	go build -o letsgo ./cmd/cli/main/main.go
	chmod +x letsgo

build-test:
	make build
	./letsgo make test test

test:
	go test -v ./internal/... 
	make build-test
	(cd test && make test)

clean:
	rm -rf ./letsgo
	rm -rf ./test
