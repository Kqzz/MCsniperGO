FROM golang:1.21.0-alpine3.18 AS build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o sniper ./cmd/cli

FROM alpine:3.18

WORKDIR /app

COPY --from=build /app/sniper .

ENTRYPOINT ["/app/sniper"]
