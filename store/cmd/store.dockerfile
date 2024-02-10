FROM alpine:latest
FROM golang:1.21.3

RUN mkdir /app
WORKDIR /app

COPY storeApp .

CMD ["./storeApp"]
