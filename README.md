# Digital Payment Platform

The **Digital Payment Platform** is a reliable, scalable, and secure integrated system designed to streamline various digital financial transactions. Built using a monolithic architecture, all core functionalities—such as card management, transactions, fund transfers, and merchant interactions—are consolidated within a single application. This approach simplifies development, testing, and deployment, enabling consistent and efficient delivery of digital payment services within a unified environment.


## 🛠️ Technologies Used
- 🚀 **gRPC** — Provides high-performance, strongly-typed APIs.
- 📡 **Kafka** — Used to publish balance-related events (e.g., after card creation).
- 📈 **Prometheus** — Collects metrics like request count and latency for each RPC method.
- 🛰️ **OpenTelemetry (OTel)** — Enables distributed tracing for observability.
- 🦫 **Go (Golang)** — Implementation language.
- 🌐 **Echo** — HTTP framework for Go.
- 🪵 **Zap Logger** — Structured logging for debugging and operations.
- 📦 **Sqlc** — SQL code generator for Go.
- 🧳 **Goose** — Database migration tool.
- 🐳 **Docker** — Containerization tool.
- 🧱 **Docker Compose** — Simplifies containerization for development and production environments.
- 🐘 **PostgreSQL** — Relational database for storing user data.
- 📃 **Swago** — API documentation generator.
- 🧭 **Zookeeper** — Distributed configuration management.
- 🔀 **Nginx** — Reverse proxy for HTTP traffic.
- 🔍 **Jaeger** — Distributed tracing for observability.
- 📊 **Grafana** — Monitoring and visualization tool.
- 🧪 **Postman** — API client for testing and debugging endpoints.
- ☸️ **Kubernetes** — Container orchestration platform for deployment, scaling, and management.
- 🧰 **Redis** — In-memory key-value store used for caching and fast data access.
- 📥 **Loki** — Log aggregation system for collecting and querying logs.
- 📤 **Promtail** — Log shipping agent that sends logs to Loki.
- 🔧 **OTel Collector** — Vendor-agnostic collector for receiving, processing, and exporting telemetry data (metrics, traces, logs).
- 🖥️ **Node Exporter** — Exposes system-level (host) metrics such as CPU, memory, disk, and network stats for Prometheus.


----

> [!WARNING]
> Important Notice: This Digital Payment Platform is currently under active development and is not production-ready. Some core features may be incomplete or subject to change. This project is intended for personal use and learning purposes only.

---

## Architecture Digital Payment Platform

### Docker

<img src="./images/archictecture_docker_payment_gateway.png" alt="docker-architecture">

### Kubernetes

<img src="./images/archictecture_kubernetes_payment_gateway.png" alt="kubernetes-architecture">



## Screenshoot

### API Documentation
<img src="./images/swagger.png" alt="hello-api-documentation">


### ERD Documentation

<img src="./images/Payment Gateway.png" alt="hello-erd-documentation" />


### Grafana Dashboard(Prometheus & OpenTelemetry(Jaeger))

#### Node Exporter

<img src="./images//node-exporter.png" alt="hello-node-exporter-grafana-dashboard">

#### Email Service

<img src="./images/email-service.png" alt="hello-email-grafana-dashboard">


#### Auth Service

<img src="./images/auth-service.png" alt="hello-auth-grafana-dashboard">

#### User Service

<img src="./images/user-service.png" alt="hello-user-grafana-dashboard">


#### Role Service

<img src="./images/role-service.png" alt="hello-role-grafana-dashboard">


#### Merchant Service

<img src="./images/merchant-service.png" alt="hello-merchant-grafana-dashboard">

#### Card Service

<img src="./images/card-service.png" alt="hello-card-grafana-dashboard">


#### Saldo Service

<img src="./images/saldo-service.png" alt="hello-saldo-grafana-dashboard">


#### Topup Service

<img src="./images/topup-service.png" alt="hello-topup-grafana-dashboard">


#### Transaction Service

<img src="./images/transaction-service.png" alt="hello-transaction-grafana-dashboard">


#### Transfer Service

<img src="./images/transfer-service.png" alt="hello-transfer-grafana-dashboard">

#### Withdraw Service

<img src="./images/withdraw-service.png" alt="hello-withdraw-grafana-dashboard">
