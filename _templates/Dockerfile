FROM golang:1.22 AS base

WORKDIR /usr/app

RUN go install github.com/air-verse/air@latest;
RUN go install github.com/a-h/templ/cmd/templ@latest;

COPY go.mod go.sum ./
RUN go mod download;
RUN go mod tidy;

COPY . .
RUN templ generate;

FROM base AS dev
RUN go build -o main cmd/http/*.go;
CMD ["air", "-c", ".air.toml"]

FROM base AS tester
CMD ["go", "test", "./..."]
