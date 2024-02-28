FROM golang:1.22-bullseye

EXPOSE 25

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o app

CMD [ "sh", "-c", "./app" ]
