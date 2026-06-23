build:
	@docker  build -t {{ .AppName }}:dev .
	@go build -o {{ .AppName }} cmd/http/main.go

release:

run: build
	@docker compose up --build


test:
	@echo -e "\n\033[0;34mBuilding image\033[m\n"
	@docker build -t {{ .AppName }}:tester --target tester . > /dev/null
	@echo -e "\n\033[0;34mRunning test image\033[m\n"
	@docker run {{ .AppName }}:tester
	@echo -e "\n\033[0;32mAll tests passed\033[m\n"

.PHONY: run build test
