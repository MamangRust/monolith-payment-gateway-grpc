# 🛡️ API Gateway Service 

A centralized entrypoint to route and manage HTTP traffic in a distributed system.
--

The **API Gateway Service** is a reverse proxy layer for a monolith-based architecture. It consolidates routing, authorization, rate-limiting, metrics collection, and service discovery for backend gRPC services such as:

- Auth Service
- User Service
- Merchant Service
- Transaction Service
- Topup, Transfer, Withdraw, etc.


## Project Structure

```
├── cmd
│   ├── main.go
│   └── README.md
├── Dockerfile
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── internal
│   ├── app
│   │   ├── client.go
│   │   └── README.md
│   ├── handler
│   │   ├── auth.go
│   │   ├── card.go
│   │   ├── handle.go
│   │   ├── merchant_document.go
│   │   ├── merchant.go
│   │   ├── README.md
│   │   ├── role.go
│   │   ├── saldo.go
│   │   ├── topup.go
│   │   ├── transaction.go
│   │   ├── transfer.go
│   │   ├── user.go
│   │   └── withdraw.go
│   └── middlewares
│       ├── auth.go
│       ├── merchant.go
│       ├── merchantRequestHandler.go
│       ├── rate_limiter.go
│       ├── README.md
│       ├── required_role.go
│       ├── role.go
│       └── roleRequestHandler.go
└── README.md
```

## HandleDomain

| Domain      | Handler File           | Endpoint Prefix          |
| ----------- | ---------------------- | ------------------------ |
| Auth        | `auth.go`              | `/api/auth`              |
| User        | `user.go`              | `/api/user`              |
| Card        | `card.go`              | `/api/card`              |
| Merchant    | `merchant.go`          | `/api/merchant`          |
| Documents   | `merchant_document.go` | `/api/merchant-document` |
| Role        | `role.go`              | `/api/role`              |
| Transaction | `transaction.go`       | `/api/transaction`                |
| Topup       | `topup.go`             | `/api/topup`             |
| Transfer    | `transfer.go`          | `/api/transfer`          |
| Withdraw    | `withdraw.go`          | `/api/withdraw`          |
| Saldo       | `saldo.go`             | `/api/saldo`             |


## 📖 API Documentation

The API docs are available via Swagger UI:

```
http://localhost:5000/swagger//
```
You can explore the available routes, request/response schemas, and try out the APIs.