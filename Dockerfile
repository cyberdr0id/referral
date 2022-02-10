FROM golang:1.16.5-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main ./cmd/main.go

ENTRYPOINT ["/app/main"]