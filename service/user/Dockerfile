FROM golang:1.23.9-alpine3.21 AS builder

RUN apk add --no-cache git

WORKDIR /app


COPY go.work go.work.sum ./
COPY pkg/ ./pkg/
COPY shared/ ./shared/
COPY service/user/ ./service/user/

WORKDIR /app/service/user

RUN go mod tidy && go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/user ./

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/user .

EXPOSE 50055 8085
CMD ["./user"]
