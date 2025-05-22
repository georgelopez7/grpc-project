# 📡 GRPC Project (+ 🔎 Observability)

This repository showcases a minimal Go-based microservices project using gRPC for communication. It simulates a simple banking backend where a payment request is processed sequentially by the following services:

- 🔍 **Validation Service** – to validate the request.
- 💀 **Fraud Service** – to check for potential fraud.

There is also a fully-implemented _Observability Stack_:

- 🔥 **Prometheus** – Metrics collection
- 🧬 **OpenTelemetry + Tempo** – Distributed tracing
- 🪓 **Loki + Promtail** - Centralized Logs
- 📊 **Grafana** – Visualization dashboard

## Getting Started

### 📂 Clone the repository

```sh
git clone https://github.com/georgelopez7/grpc-project.git
cd grpc-project
```

### 🐋 Start the project

**_NOTE:_** Ensure you have `Docker` running

```sh
docker compose up --build
```

This will start up the following services:

#### 📦 Microservices:

- **gateway** - HTTP server running on port `:8080`
- **fraud service** - gRPC server running on port `:50052`
- **validation service** - gRPC server running on port `:50053`

#### 🔎 Observability Services:

- **Grafana** - running at `http://localhost:3000`
- **Prometheus (metrics)**
- **Tempo (traces)**
- **Loki + Promtail (logs)**

### 🛰 Send HTTP Request

Send `POST` request to `http://localhost:8080/api/v1/payments`

**Request Body:**

```json
{
  "id": "hfuwahf",
  "sender": "george",
  "receiver": "lopez",
  "amount": 100
}
```

**Response:**

```json
{
  "fraud_status": {
    "is_fraudulent": false,
    "message": "Payment is not fraudulent!"
  },
  "validation_status": {
    "is_valid": true,
    "message": "Payment passed validation check!"
  }
}
```

The response tells us if the payment is both **valid** & if the payment is **fraudulent.**

## 📚 Documentation

### 🏗 Project Architecture

Below is a diagram of the flow of a request through the system:

![grpc-project-diagram](https://github.com/user-attachments/assets/0da5ee8a-689b-42cb-81a8-12d4dbaca3f5)

---

### 📦 Microservices

#### 🚧 Gateway

The **gateway** is the entry point for all payment requests. It receives **HTTP requests** and routes the payment data to the **Validation** and **Fraud** microservices for processing.

Built using the Gin web framework.

**API Endpoint**

**POST** `/api/v1/payments` --> Submits a payment request.

**_Response:_** Indicates whether the request is **valid** and **fraud-free.**

---

#### 🔍 Validation Service

The **Validation Service** checks if a payment request is _valid_ based on a simple rule:

- The payment amount must be less than **1000** _(default limit)._

This service runs as a **gRPC server**, receives the payment data from the **gateway**, and returns a validation result.

---

#### 💀 Fraud Service

The **Fraud Service** determines if a payment request is _fraudulent._

- A payment request is flagged as **fraudulent** if the payment amount is a **Fibonacci number** (e.g. 2, 3, 5, 8, etc.).

This service runs as a **gRPC server**, receives payment data from the **Gateway**, and returns a fraud check result.

---

### 🔍 Observability Services

#### 🔥 Prometheus _(Metrics)_

We use **Prometheus** in this project to gather metrics about our services.

**Custom Metrics**

- `gateway_payment_requests_total` - tracks the number of requests sent to the **gateway** service in total overtime.

**Config**

Location: `pkg/logging/config/prometheus-config.yaml`

**Docker**

Below is how we define `prometheus` in the `docker-compose.yml`:

```bash
prometheus:
    image: prom/prometheus
    volumes:
      - ./pkg/logging/config/prometheus-config.yaml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    ports:
      - "9090:9090" # PROMETHEUS PORT

```

---

#### 🧬 Opentelemetry & Tempo _(Traces)_

We use **OpenTelemetry** and **Tempo** for **distributed tracing.**
This setup allows you to trace the full lifecycle of a request across services.

**Config**

Location: `pkg/logging/config/tempo-config.yaml`

**Docker**

Below is how we define `tempo` in the `docker-compose.yml`:

```bash
tempo:
    image: grafana/tempo
    command: ["-config.file=/etc/tempo/tempo.yaml"]
    volumes:
        - ./pkg/logging/config/tempo-config.yaml:/etc/tempo/tempo.yaml
    ports:
        - "4318:4318" # OPENTELEMETRY PORT
        - "3200:3200" # TEMPO PORT

```

---

#### 🪓 Loki & Promtail _(Logs)_

We use **Loki** and **Promtail** for **centralized logs.**
This is so we can store and visualize logs from all different services in one place.

_Loki_ - stores the logs & in used by Grafana for visualization

_Promtail_ - scrapes logs from the Docker container and pushes data to Loki

**Config(s)**

_Loki_ - `pkg/logging/config/loki-config.yaml`

_Promtail_ - `pkg/logging/config/promtrail-config.yaml`

**Docker**

Below is how we define `loki` and `promtail` in the `docker-compose.yml`:

```bash
loki:
    image: grafana/loki:2.9.5
    ports:
      - "3100:3100"
    volumes:
      - ./pkg/logging/config/loki-config.yaml:/etc/loki/config.yaml
      - loki_data:/loki
    command: -config.file=/etc/loki/config.yaml

promtail:
    image: grafana/promtail:2.9.5
    volumes:
        - ./pkg/logging/config/promtail-config.yaml:/etc/promtail/config.yml
        - /var/run/docker.sock:/var/run/docker.sock:ro
        - promtail_positions:/tmp
    command: -config.file=/etc/promtail/config.yml
    depends_on:
        - loki
```

---

#### 📊 Grafana _(Visualization dashboards)_

We use **Grafana** to visualize all the data collected from the services _(metrics, traces & logs)_

**Config**

We have a unique folder for _Grafana_ located at `/pkg/logging/grafana`. This allows instance connection to _Prometheus_, _Tempo_ & _Loki_.

**Docker**

Below is how we define `grafana` in the `docker-compose.yml`:

```bash
grafana:
    image: grafana/grafana
    ports:
        - "3000:3000"
    volumes:
        - grafana_data:/var/lib/grafana
        - ./pkg/logging/grafana:/etc/grafana/provisioning/datasources
    depends_on:
        - prometheus
        - tempo
        - loki
    environment:
        GF_AUTH_ANONYMOUS_ENABLED: "true" # Optional: For easier dev access without login
        GF_AUTH_ANONYMOUS_ORG_ROLE: Admin # Optional: Give anonymous users admin rights for dev
        GF_AUTH_DISABLE_SIGNOUT_MENU: "true" # Optional: Remove sign-out button
        GF_AUTH_DISABLE_LOGIN_FORM: "true" # Optional: Disable login form if anonymous is enabled
        GF_THEME: light # Set default theme to light

```

## 📝 Extra Tips

#### 😭 Too Much Data

You can clear all the **volumes** in the `docker-compose.yml` by running:

```bash
docker-compose down -v
```

#### 💪🏻 Updating Protobuf

If you wish to **alter** the **protobuf schema** you can use the following command found in the `Makefile`:

```bash
make genproto-payment
```
