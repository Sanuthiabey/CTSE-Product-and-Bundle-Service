# Quick Start Guide

## Prerequisites
- Go 1.21+
- PostgreSQL (or Neon Database)
- Protocol Buffers compiler (protoc)

## Setup

1. **Set Environment Variable**
```bash
# Windows PowerShell
$env:DATABASE_URL="your-postgresql-connection-string"

# Example:
$env:DATABASE_URL="postgresql://user:password@host/database?sslmode=require"
```

2. **Install Dependencies**
```bash
go mod download
```

3. **Build the Application**
```bash
go build -o bin/server.exe ./cmd/server
```

4. **Run the Server**
```bash
.\bin\server.exe
```

The server will start:
- **REST API** on port `8080`
- **gRPC Server** on port `50053`

## Quick Test

### 1. Create a Product
```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "id": "prod-001",
    "name": "Lavender Essential Oil",
    "description": "Calming lavender oil",
    "price": 19.99,
    "mood": "calm",
    "category": "essential-oils",
    "image": "lavender.jpg",
    "rating": 4.5,
    "featured": true,
    "stock": 100
  }'
```

### 2. Create Another Product
```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "id": "prod-002",
    "name": "Chamomile Tea",
    "description": "Soothing chamomile tea",
    "price": 9.99,
    "mood": "calm",
    "category": "beverages",
    "image": "chamomile.jpg",
    "rating": 4.7,
    "featured": false,
    "stock": 50
  }'
```

### 3. Create a Bundle
```bash
curl -X POST http://localhost:8080/bundles \
  -H "Content-Type: application/json" \
  -d '{
    "id": "bundle-001",
    "name": "Relaxation Bundle",
    "mood": "calm",
    "products": [
      {"product_id": "prod-001", "quantity": 2},
      {"product_id": "prod-002", "quantity": 1}
    ]
  }'
```

### 4. Validate Stock
```bash
curl -X POST http://localhost:8080/stock/validate \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "prod-001",
    "quantity": 5
  }'
```

**Expected Response:**
```json
{
  "available": true,
  "current_stock": 100,
  "message": "Stock available"
}
```

### 5. Reduce Stock
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

**Expected Response:**
```json
{
  "success": true,
  "message": "Stock reduced successfully"
}
```

### 6. Get All Bundles
```bash
curl http://localhost:8080/bundles
```

### 7. Get Bundle by ID
```bash
curl http://localhost:8080/bundles/bundle-001
```

### 8. Update Bundle
```bash
curl -X PUT http://localhost:8080/bundles/bundle-001 \
  -H "Content-Type: application/json" \
  -d '{
    "id": "bundle-001",
    "name": "Premium Relaxation Bundle",
    "mood": "calm",
    "products": [
      {"product_id": "prod-001", "quantity": 3}
    ]
  }'
```

### 9. Delete Bundle
```bash
curl -X DELETE http://localhost:8080/bundles/bundle-001
```

## Development

### Regenerate Protobuf Files
If you modify `proto/product.proto`, regenerate the Go code:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/product.proto
```

### Run with Docker
```bash
docker-compose up --build
```

## API Endpoints Summary

### Products
- `POST /products` - Create product
- `GET /products` - List products (supports filters: mood, category, search, featured, sort)
- `GET /products/:id` - Get product by ID
- `PUT /products/:id` - Update product
- `DELETE /products/:id` - Delete product

### Bundles
- `POST /bundles` - Create bundle
- `GET /bundles` - List all bundles
- `GET /bundles/:id` - Get bundle by ID
- `PUT /bundles/:id` - Update bundle
- `DELETE /bundles/:id` - Delete bundle

### Stock Operations
- `POST /stock/validate` - Validate stock availability
- `POST /stock/reduce` - Reduce stock (transactional)

### gRPC Services
- `ValidateBundle(BundleRequest)` - Validate bundle stock
- `DeductBundle(BundleRequest)` - Deduct bundle stock
- `ValidateStock(StockRequest)` - Validate product stock
- `ReduceStock(StockReductionRequest)` - Reduce product stock

## Health Check
```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "service": "product-and-bundle-service",
  "status": "running"
}
```

## Troubleshooting

### Database Connection Issues
- Verify `DATABASE_URL` is set correctly
- Check network connectivity to database
- Ensure database user has proper permissions

### Port Already in Use
If port 8080 or 50053 is already in use:
- Stop the conflicting service
- Or modify the ports in `cmd/server/main.go` and `internal/grpc/server.go`

### CORS Issues
The service allows requests from `http://localhost:3000` by default.
To modify, update the CORS configuration in `cmd/server/main.go`.

## Next Steps

See [API_DOCUMENTATION.md](./API_DOCUMENTATION.md) for complete API reference with examples.

