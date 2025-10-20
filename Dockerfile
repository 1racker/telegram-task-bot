FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o bot main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/bot /bot
EXPOSE 8080
CMD ["./bot"]