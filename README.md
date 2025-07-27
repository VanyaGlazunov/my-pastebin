# MyPastebin

This project is my first step on the path of becoming observability monster. MyPastebin is a fully-featured, observable, and load-testable implementation of a minimalist Pastebin-like service. 

The core functionality allows users to anonymously share text snippets (pastes) with a defined expiration time, after which they are automatically deleted.

## 🚀 Core Concepts & Tech Stack

- **Backend**: Go (`gin` framework)
- **Database**: PostgreSQL (`gorm` ORM)
- **API Documentation**: Swagger (OpenAPI)
- **Containerization**: Docker & Docker Compose
- **Monitoring & Metrics**: Prometheus
- **Visualization & Dashboards**: Grafana
- **Alerting**: Alertmanager with Telegram notifications
- **Distributed Tracing**: Jaeger with OpenTelemetry
- **Load Testing**: Locust

## ✨ Features

- **Create Pastes**: Submit text content with a specified expiration duration (e.g., `10m`, `1h`, `24h`).
- **View Pastes**: Retrieve saved pastes via a unique short ID.
- **Automatic Deletion**: A background worker periodically cleans up expired pastes from the database.
- **Full Observability**: Pre-configured dashboards in Grafana for application and database performance.
- **Actionable Alerts**: Proactive notifications sent to a Telegram channel for critical performance issues.
- **Distributed Tracing**: End-to-end request tracing from the load tester to the database.
- **Interactive Load Testing**: A web UI to control load profiles and stress-test the system.
- **API Documentation**: Interactive Swagger UI for easy API exploration.

---

## Accessing Services
I have encountered eventual problems in firefox, always worked in yandex-browser though

- **Docs & Links Hub**: http://51.250.3.188:8080/docs
- **API (Swagger UI)**: http://51.250.3.188:8080/swagger/index.html
- **Monitoring (Grafana)**: http://51.250.3.188:3000 (Login: admin / admin)
- **Tracing (Jaeger)**: http://51.250.3.188:16686
- **Load Testing (Locust)**: http://51.250.3.188:8089
- **Alerts (Alertmanager)**: http://51.250.3.188:9093
- **Metrics (Prometheus)**: http://51.250.3.188:9090
- **Telegram Channel For Alerts** https://t.me/MyPastebinAlerts

## 🚨 Triggering Alerts
1. Start test in Locust for 
2. Wait for a few minutes, Observe p99 and DBRPS graphs hits red zones.
3. View alerts in telegram channel.

