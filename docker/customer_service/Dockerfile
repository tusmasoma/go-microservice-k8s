FROM --platform=linux/amd64 golang:1.23.4 AS builder

RUN apt-get update && apt-get install -y default-mysql-client

WORKDIR /app

COPY ./services/customer/go.mod ./
COPY ./services/customer/go.sum ./

RUN go mod download

COPY ./services/customer/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

FROM --platform=linux/amd64 alpine:3.18 AS production

RUN apk add --no-cache mysql-client

WORKDIR /app

COPY --from=builder /app/main .

COPY ./hack/entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["entrypoint.sh"]

CMD ["/app/main"]
