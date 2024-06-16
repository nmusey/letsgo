FROM golang:1.22 as base

WORKDIR /usr/app

RUN go install github.com/air-verse/air@latest;
RUN go install github.com/a-h/templ/cmd/templ@latest;

COPY go.mod go.sum ./
RUN go mod download;
RUN go mod verify;
RUN go mod tidy;

COPY . .
RUN templ generate;

FROM base as runner
RUN go build -o main cmd/http/*.go;
CMD ["air", "-c", ".air.toml"]

FROM base as tester
CMD ["go", "test", "./..."]
