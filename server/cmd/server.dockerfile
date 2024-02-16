FROM alpine:latest
FROM golang:1.21.3

RUN mkdir /app

WORKDIR /app

COPY serverApp .
EXPOSE 8080
CMD ["./serverApp"]
