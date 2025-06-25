FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /bin/app ./cmd/server

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /bin/app .
COPY .env ./
COPY exports ./exports

EXPOSE 8002

CMD ["./app"]
