# syntax=docker/dockerfile:1
FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o seaport main.go

FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/seaport /app/seaport
EXPOSE 8080
ENTRYPOINT ["/app/seaport"]
