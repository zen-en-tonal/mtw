FROM golang:1.22-alpine3.19 AS builder

WORKDIR /app
COPY . .

RUN apk update && apk add alpine-sdk
RUN go mod download
RUN go build cmd/server/serve.go

FROM alpine:3.19

EXPOSE 25
EXPOSE 8080

ENV GIN_MODE=release

COPY migrations migrations
RUN mkdir db
COPY --from=builder /app/serve .

CMD ["./serve"]
