FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o cloudrun ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/cloudrun .
COPY --from=builder /app/pkg/weather_api/.env ./pkg/weather_api/.env
ENTRYPOINT ["/app/cloudrun"]
