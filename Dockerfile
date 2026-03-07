FROM golang:1.26-alpine3.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main ./cmd/api

FROM alpine:3.23
WORKDIR /app
COPY --from=builder /app/main main
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["./main"]
