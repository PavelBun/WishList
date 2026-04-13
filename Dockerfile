FROM golang:1.25.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /wishlist-api ./cmd/server

FROM alpine:3.18
RUN adduser -D -u 1000 appuser
WORKDIR /app
COPY --from=builder /wishlist-api /app/
COPY migrations /app/migrations
USER appuser
EXPOSE 8080
CMD ["./wishlist-api"]