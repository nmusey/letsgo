run: build
	@docker compose up --build

build:
	@docker  build -t {{ .AppName }}:dev .
	@go build -o {{ .Substitutions.AppName }} cmd/http/main.go

test:
	@echo -e "\n\033[0;34mBuilding image\033[m\n"
	@docker build -t {{ .Substitutions.AppName }}:tester --target tester . > /dev/null
	@echo -e "\n\033[0;34mRunning test image\033[m\n"
	@docker run {{ .Substitutions.AppName }}:tester
	@echo -e "\n\033[0;32mAll tests passed\033[m\n"

.PHONY: run build test
