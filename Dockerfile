FROM golang:1.19.4 AS builder
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./pkg ./pkg

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./build/main ./cmd/lambda/main.go

FROM public.ecr.aws/lambda/go:1
COPY --from=builder app/build/main ${LAMBDA_TASK_ROOT}
CMD [ "main" ]