FROM golang:1.22

WORKDIR /usr/app

RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/a-h/templ/cmd/templ@latest

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go mod tidy
RUN templ generate
RUN go build -o main cmd/http/*.go

CMD ["air", "-c", ".air.toml"]
