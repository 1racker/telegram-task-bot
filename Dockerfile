FROM golang:1.25.0-alpine AS builder
RUN apk add --no-cache gcc g++ musl-dev postgresql-dev
ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o bot main.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates libpq tzdata
COPY --from=builder /app/bot /usr/local/bin/bot
RUN adduser -D appuser
USER appuser
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/bot"]