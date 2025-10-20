FROM golang:1.25.0-alpine AS builder
RUN apk add --no-cache gcc g++ musl-dev sqlite-dev
ENV CGO_ENABLED=1
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o bot main.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates sqlite-libs
COPY --from=builder /app/bot /bot
EXPOSE 8080
CMD ["./bot"]