FROM golang:1.24-alpine AS build

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY configs ./configs
COPY internal ./internal
COPY .env ./.env

RUN GOOS=linux GOARCH=amd64 go build -o ./bin/main ./cmd

FROM alpine:latest AS app

WORKDIR /app

EXPOSE 8080

COPY --from=build /app/bin/main .

CMD ["./main"]
