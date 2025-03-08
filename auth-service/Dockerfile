FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY --from=builder /app/cmd/migrate/migrations ./cmd/migrate/migrations
COPY --from=builder /app/private_key_access.pem ./private_key_access.pem
COPY --from=builder /app/private_key_refresh.pem ./private_key_refresh.pem
EXPOSE 8080
CMD ["./main"]