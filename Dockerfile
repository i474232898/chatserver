# syntax=docker/dockerfile:1

FROM golang:1.24.4 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /chatserver ./cmd/chatserver/main.go

FROM alpine:3.19.0 AS runtime-stage

WORKDIR /app

COPY --from=build-stage /chatserver /app/chatserver
COPY --from=build-stage /app/api /app/api

EXPOSE 8080

ENTRYPOINT ["/app/chatserver"]
