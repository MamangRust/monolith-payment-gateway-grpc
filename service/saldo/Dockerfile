FROM golang:1.23.9-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o saldo ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
RUN mkdir -p /var/log/app

COPY --from=builder /app/saldo .

EXPOSE 50056 8086
CMD ["./saldo"]