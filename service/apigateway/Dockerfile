FROM golang:1.23.9-alpine3.21 AS builder

WORKDIR /app

COPY service/apigateway/go.mod service/apigateway/go.sum ./
RUN go mod download


COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o apigateway ./service/apigateway/cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/apigateway .

EXPOSE 5000 50051
CMD ["./apigateway"]
