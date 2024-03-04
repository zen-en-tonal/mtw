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

CMD ["./serve"]
