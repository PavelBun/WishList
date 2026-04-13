FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /wishlist-api ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /wishlist-api .

# По умолчанию порт 8080, но можно переопределить через ENV
ENV APP_PORT=8080
EXPOSE ${APP_PORT}

CMD ["sh", "-c", "./wishlist-api"]