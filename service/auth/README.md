# 🛡️ Auth Service

Welcome to the **Auth Service** — a modular and extensible service that provides authentication, authorization, user registration, password reset, and identity verification functionalities. This service is designed to be part of a microservice architecture and can be deployed independently.

---

## 🧭 Project Structure

```bash
├── cmd/                     # App entrypoint
│   └── main.go
├── internal/
│   ├── apps/                # HTTP server setup and dependencies
│   ├── errorhandler/        # Centralized error handling
│   ├── handler/             # Route and request handlers
│   ├── redis/               # Redis caching logic per use-case
│   ├── repository/          # Interfaces + persistence layer (e.g., RefreshToken, User)
│   └── service/             # Business logic for Auth, Register, Identity, etc.
├── Dockerfile               # Container definition
├── go.mod / go.sum          # Go module files
└── README.md                # This file
```


