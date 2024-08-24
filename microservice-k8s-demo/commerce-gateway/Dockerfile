FROM golang:1.21.3

RUN apt-get update && apt-get install -y default-mysql-client

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

CMD ["/app/main"]
