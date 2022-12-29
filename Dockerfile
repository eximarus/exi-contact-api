FROM golang:1.19.4
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./pkg ./pkg

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./build/main ./cmd/local/main.go

CMD [ "build/main" ]