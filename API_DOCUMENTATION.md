# API Documentation

## Overview
This service provides REST and gRPC APIs for managing products, bundles, and stock operations.

**Base URL (REST)**: `http://localhost:8080`  
**gRPC Port**: `50053`

---

## 1️⃣ Bundle Management (CRUD)

### Create Bundle
**Endpoint**: `POST /bundles`

**Request Body**:
```json
{
  "id": "bundle-001",
  "name": "Relaxation Bundle",
  "mood": "calm",
  "products": [
    {
      "product_id": "prod-001",
      "quantity": 2
    },
    {
      "product_id": "prod-002",
      "quantity": 1
    }
  ]
}
```

**Response**: `201 Created`
```json
{
  "id": "bundle-001",
  "name": "Relaxation Bundle",
  "mood": "calm",
  "products": [
    {
      "product_id": "prod-001",
      "quantity": 2
    },
    {
      "product_id": "prod-002",
      "quantity": 1
    }
  ]
}
```

---

### Get All Bundles
**Endpoint**: `GET /bundles`

**Response**: `200 OK`
```json
[
  {
    "id": "bundle-001",
    "name": "Relaxation Bundle",
    "mood": "calm",
    "products": [
      {
        "product_id": "prod-001",
        "quantity": 2
      },
      {
        "product_id": "prod-002",
        "quantity": 1
      }
    ]
  }
]
```

---

### Get Bundle by ID
**Endpoint**: `GET /bundles/:id`

**Example**: `GET /bundles/bundle-001`

**Response**: `200 OK`
```json
{
  "id": "bundle-001",
  "name": "Relaxation Bundle",
  "mood": "calm",
  "products": [
    {
      "product_id": "prod-001",
      "quantity": 2
    },
    {
      "product_id": "prod-002",
      "quantity": 1
    }
  ]
}
```

**Error Response**: `404 Not Found`
```json
{
  "error": "Bundle not found"
}
```

---

### Update Bundle
**Endpoint**: `PUT /bundles/:id`

**Example**: `PUT /bundles/bundle-001`

**Request Body**:
```json
{
  "id": "bundle-001",
  "name": "Updated Relaxation Bundle",
  "mood": "calm",
  "products": [
    {
      "product_id": "prod-001",
      "quantity": 3
    }
  ]
}
```

**Response**: `200 OK`
```json
{
  "id": "bundle-001",
  "name": "Updated Relaxation Bundle",
  "mood": "calm",
  "products": [
    {
      "product_id": "prod-001",
      "quantity": 3
    }
  ]
}
```

---

### Delete Bundle
**Endpoint**: `DELETE /bundles/:id`

**Example**: `DELETE /bundles/bundle-001`

**Response**: `200 OK`
```json
{
  "message": "Bundle deleted"
}
```

**Error Response**: `404 Not Found`
```json
{
  "error": "Bundle not found"
}
```

---

## 2️⃣ Stock Validation API

### REST API: Validate Stock
**Endpoint**: `POST /stock/validate`

**Request Body**:
```json
{
  "product_id": "prod-001",
  "quantity": 5
}
```

**Response**: `200 OK` (Stock Available)
```json
{
  "available": true,
  "current_stock": 10,
  "message": "Stock available"
}
```

**Response**: `200 OK` (Insufficient Stock)
```json
{
  "available": false,
  "current_stock": 3,
  "message": "Insufficient stock"
}
```

**Response**: `404 Not Found` (Product Not Found)
```json
{
  "available": false,
  "current_stock": 0,
  "message": "Product not found"
}
```

---

### gRPC: ValidateStock
**Service**: `ProductService`  
**Method**: `ValidateStock`

**Request**:
```protobuf
message StockRequest {
  string product_id = 1;
  int32 quantity = 2;
}
```

**Response**:
```protobuf
message StockValidateResponse {
  bool available = 1;
  int32 current_stock = 2;
  string message = 3;
}
```

**Example (Go Client)**:
```go
req := &pb.StockRequest{
    ProductId: "prod-001",
    Quantity: 5,
}

res, err := client.ValidateStock(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Available: %v, Stock: %d, Message: %s\n", 
    res.Available, res.CurrentStock, res.Message)
```

---

## 3️⃣ Stock Reduction API

### REST API: Reduce Stock
**Endpoint**: `POST /stock/reduce`

**Request Body**:
```json
{
  "items": [
    {
      "product_id": "prod-001",
      "quantity": 2
    },
    {
      "product_id": "prod-002",
      "quantity": 1
    }
  ]
}
```

**Response**: `200 OK` (Success)
```json
{
  "success": true,
  "message": "Stock reduced successfully"
}
```

**Response**: `400 Bad Request` (Insufficient Stock)
```json
{
  "success": false,
  "message": "Insufficient stock for product prod-001"
}
```

**Response**: `404 Not Found` (Product Not Found)
```json
{
  "success": false,
  "message": "Product prod-001 not found"
}
```

---

### gRPC: ReduceStock
**Service**: `ProductService`  
**Method**: `ReduceStock`

**Request**:
```protobuf
message StockReductionRequest {
  repeated StockItem items = 1;
}

message StockItem {
  string product_id = 1;
  int32 quantity = 2;
}
```

**Response**:
```protobuf
message StockReductionResponse {
  bool success = 1;
  string message = 2;
}
```

**Example (Go Client)**:
```go
req := &pb.StockReductionRequest{
    Items: []*pb.StockItem{
        {ProductId: "prod-001", Quantity: 2},
        {ProductId: "prod-002", Quantity: 1},
    },
}

res, err := client.ReduceStock(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Success: %v, Message: %s\n", res.Success, res.Message)
```

---

## Existing gRPC APIs

### ValidateBundle
**Method**: `ValidateBundle`

Validates if all products in a bundle have sufficient stock.

**Request**:
```protobuf
message BundleRequest {
  string bundle_id = 1;
}
```

**Response**:
```protobuf
message ValidateResponse {
  bool valid = 1;
  string message = 2;
}
```

---

### DeductBundle
**Method**: `DeductBundle`

Deducts stock for all products in a bundle (transactional).

**Request**:
```protobuf
message BundleRequest {
  string bundle_id = 1;
}
```

**Response**:
```protobuf
message DeductResponse {
  bool success = 1;
  string message = 2;
}
```

---

## Testing Examples

### cURL Examples

#### Create a Bundle
```bash
curl -X POST http://localhost:8080/bundles \
  -H "Content-Type: application/json" \
  -d '{
    "id": "bundle-001",
    "name": "Wellness Bundle",
    "mood": "energetic",
    "products": [
      {"product_id": "prod-001", "quantity": 2},
      {"product_id": "prod-002", "quantity": 1}
    ]
  }'
```

#### Validate Stock
```bash
curl -X POST http://localhost:8080/stock/validate \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "prod-001",
    "quantity": 5
  }'
```

#### Reduce Stock
```bash
curl -X POST http://localhost:8080/stock/reduce \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {"product_id": "prod-001", "quantity": 2},
      {"product_id": "prod-002", "quantity": 1}
    ]
  }'
```

---

## Error Handling

All endpoints follow standard HTTP status codes:
- `200 OK`: Successful operation
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request format or insufficient stock
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

Error responses follow this format:
```json
{
  "error": "Error message description"
}
```

---

## Database Schema

### bundles
```sql
CREATE TABLE bundles (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  mood TEXT
);
```

### bundle_products
```sql
CREATE TABLE bundle_products (
  bundle_id TEXT REFERENCES bundles(id) ON DELETE CASCADE,
  product_id TEXT REFERENCES products(id) ON DELETE CASCADE,
  quantity INT NOT NULL,
  PRIMARY KEY (bundle_id, product_id)
);
```

### products
```sql
CREATE TABLE products (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  price DOUBLE PRECISION,
  mood TEXT,
  category TEXT,
  image TEXT,
  rating DOUBLE PRECISION DEFAULT 0,
  featured BOOLEAN DEFAULT FALSE,
  stock INT
);
```

---

## Notes

1. **Transaction Safety**: Both Stock Reduction APIs use database transactions to ensure atomicity.
2. **Cascade Deletes**: Deleting a bundle automatically removes its associated products from `bundle_products`.
3. **Validation**: All stock operations validate availability before deduction.
4. **CORS**: The service is configured to accept requests from `http://localhost:3000`.


