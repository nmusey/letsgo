FROM golang:1.21

WORKDIR /usr/app

RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o main .

CMD ["air", "-c", ".air.toml"]