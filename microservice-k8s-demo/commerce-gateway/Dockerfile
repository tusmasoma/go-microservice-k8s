FROM golang:1.21.3 AS builder

RUN apt-get update && apt-get install -y default-mysql-client

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

FROM alpine:3.18 AS production

RUN apk add --no-cache mysql-client

WORKDIR /app

COPY --from=builder /app/main .

CMD ["/app/main"]
