FROM golang:1.24-bookworm

WORKDIR /app

COPY . .
RUN go mod tidy

RUN go build -o kv-server ./cmd/server

EXPOSE 8080

CMD ["./kv-server"]
