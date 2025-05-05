# gRPC Project

This repository showcases a minimal Go-based microservices project using gRPC for communication. It simulates a simple banking backend where a payment request is processed sequentially by the following services:

- Validation Service – to validate the request.
- Fraud Detection Service – to check for potential fraud.

Observability Stack:

- Prometheus – Metrics collection
- OpenTelemetry + Tempo – Distributed tracing
- Grafana – Visualization dashboard

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

- gateway --> HTTP server running on port `:8080`
- fraud service --> gRPC server running on port `:50052`
- validation service --> gRPC server running on port `:50053`

#### 🔍 Observability Services:

- Grafana --> running at `localhost:3000`
- Prometheus (metrics)
- Tempo (traces)

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

The response tells us if the payment is both valid & if the payment is fraudulent.

## 📚 Documentation

### 🏗 Project Architecture

Below is a diagram of the flow of a request through the system:

![grpc-project-diagram](https://github.com/user-attachments/assets/0da5ee8a-689b-42cb-81a8-12d4dbaca3f5)

---

### 📦 Microservices

#### 🚧 Gateway

The **Gateway** is the entry point for all payment requests. It receives HTTP requests and routes the payment data to the **Validation** and **Fraud Detection** microservices for processing.

Built using the Gin web framework.

**API Endpoint**

**POST** `/api/v1/payments` --> Submits a payment request.

**_Response:_** Indicates whether the request is valid and fraud-free.

---

#### 🔍 Validation Service

The **Validation Service** checks if a payment request is _valid_ based on a simple rule:

- The payment amount must be less than **1000** (default limit).

This service runs as a **gRPC server**, receives the payment data from the **Gateway**, and returns a validation result.

---

#### ❌ Fraud Detection Service

The **Fraud Service** determines if a payment request is _fraudulent._

A request is flagged as **fraudulent** if the payment amount is a **Fibonacci number** (e.g. 2, 3, 5, 8, etc.).

This service runs as a **gRPC server**, receives payment data from the **Gateway**, and returns a fraud check result.

---

### 🔍 Observability Services

#### 🔥 Prometheus _(Metrics)_

We use Prometheus in this project to gather metrics about our services.

**Manually Recorded Metrics**

- `gateway_payment_requests_total` --> tracks the number of requests sent to the **Gateway** service in total overtime.

---

#### 💨 Opentelemetry & Tempo _(Traces)_

**_NOTE:_** Traces have only been set up for the `gateway` service at this time.

The project uses **OpenTelemetry** and **Tempo** for distributed tracing.
This setup allows you to trace the full lifecycle of a request across services and quickly identify issues or bottlenecks.

---

## 🗻 Next Steps

Below is a list of improvements that I wish to add:

- Cross service traces --> pass a trace from gateway to fraud service for example
- Logs using Loki and displayed using Grafana
- _Coming soon..._
