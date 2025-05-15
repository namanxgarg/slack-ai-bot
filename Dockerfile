FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o slack-bot ./cmd/main.go

EXPOSE 3000
CMD ["./slack-bot"]
