# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/migrate cmd/migrate/main.go

# Final Stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/api /app/api
COPY --from=builder /app/bin/migrate /app/migrate

# Add non-root user for security
RUN adduser -D appuser
USER appuser

EXPOSE 8080

CMD ["./api"]
