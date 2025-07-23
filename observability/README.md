# 📡 Observability Configuration

This directory contains the **observability stack configuration** for the microservices platform. It includes definitions for monitoring, logging, tracing, and alerting components using tools like **Prometheus**, **Alertmanager**, **Grafana**, **Loki**, **Promtail**, and **OpenTelemetry Collector**.


## Project Structure

```
observability/
├── alertmanager.yml # Alertmanag er routing and notification configuration
├── loki-config.yaml # Loki logging backend configuration 
├── otel-collector.yaml  # OpenTelemetry Collector pipelines and exporters
├── prometheus.yaml  # Prometheus scrape configs and alert rule loading
├── promtail-config.yaml # Promtail log shipper configuration for Loki
├── README.md
└── rules  # Directory for service-specific alert rules
    ├── apigateway-alerts.yaml
    ├── auth-service-alerts.yaml
    ├── card-service-alerts.yaml
    ├── email-service-alerts.yaml
    ├── golang-runtime-alerts.yaml
    ├── kafka-exporter-alerts.yaml
    ├── merchant-service-alerts.yaml
    ├── node-exporter-alerts.yaml
    ├── otel-collector-alerts.yaml
    ├── role-service-alerts.yaml
    ├── saldo-service-alerts.yaml
    ├── topup-service-alerts.yaml
    ├── transaction-service-alerts.yaml
    ├── transfer-service-alerts.yaml
    ├── user-service-alerts.yaml
    └── withdrawal-service-alerts.yaml
```