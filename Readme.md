# 📦Product & Bundle Service

A Go-based microservice for managing products, bundles, and stock operations.

Built using:

- Go (Gin)
- PostgreSQL
- Docker
- gRPC
- REST API

---

## 🚀 Features

### 🛍 Product Management (REST)
- Create product
- Get all products
- Get product by ID
- Update product
- Delete product

### 📦 Bundle Management (REST)
- Create bundle
- Get all bundles
- Get bundle by ID

### 🔄 Stock Operations
- Validate bundle stock
- Deduct bundle stock (transactional)

### ⚡ Internal Communication
- gRPC server for microservice communication
- Designed to be called by Order Service

---

## 🏗 Architecture

This service is part of a microservices architecture:

Client → Order Service → Product Service  
↳ Auth Service

- Order Service calls Product Service via gRPC
- Product Service handles stock validation & deduction
- Authentication handled separately

---

## 🛠 Tech Stack

- Go 1.25
- Gin (REST API)
- PostgreSQL 15
- gRPC
- Docker & Docker Compose

---

### 🌐 API Runs On
http://localhost:8080

### ⚡ gRPC Runs On
localhost:50051

### 🔐 gRPC Methods

| Method          | Description                       |
|-----------------|-----------------------------------|
| ValidateBundle  | Checks stock availability         |
| DeductBundle    | Deducts stock using transaction   |

## 📂 Project Structure

- cmd/server/ → Main application entry
- internal/db/ → Database connection
- internal/grpc/ → gRPC server implementation
- internal/models/ → Data models
- proto/ → gRPC protobuf definitions
- docker-compose.yml → Container orchestration

---

## ▶️ Running the Service

### Build & Run with Docker

```bash
docker compose down
docker compose up --build