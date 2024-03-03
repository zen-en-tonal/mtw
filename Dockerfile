FROM golang:1.22-alpine3.19 AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build cmd/server/serve.go

FROM alpine:3.19

EXPOSE 25
EXPOSE 8080

ENV GIN_MODE=release

WORKDIR /app
COPY migrations migrations
COPY --from=builder /app/serve .

RUN apk add --no-cache ca-certificates && update-ca-certificates
ENV SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt
ENV SSL_CERT_DIR=/etc/ssl/certs

CMD ["./serve"]
