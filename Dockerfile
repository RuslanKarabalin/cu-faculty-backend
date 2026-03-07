FROM golang:1.26-alpine3.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main ./cmd/api

FROM alpine:3.23
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /app/main main
COPY --from=builder /app/static static
EXPOSE 8080
USER appuser:appgroup
ENTRYPOINT ["./main"]
