package grpc

import (
	"context"
	"log"
	"net"

	"github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/internal/db"
	pb "github.com/Sanuthiabey/CTSE-Product-and-Bundle-Service/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedProductServiceServer
}

// -----------------------------
// VALIDATE BUNDLE
// -----------------------------
func (s *server) ValidateBundle(ctx context.Context, req *pb.BundleRequest) (*pb.ValidateResponse, error) {

	query := `
	SELECT p.stock, bp.quantity
	FROM bundle_products bp
	JOIN products p ON p.id = bp.product_id
	WHERE bp.bundle_id = $1;
	`

	rows, err := db.DB.Query(query, req.BundleId)
	if err != nil {
		return &pb.ValidateResponse{Valid: false, Message: err.Error()}, nil
	}
	defer rows.Close()

	found := false

	for rows.Next() {
		found = true
		var stock int
		var quantity int

		if err := rows.Scan(&stock, &quantity); err != nil {
			return &pb.ValidateResponse{Valid: false, Message: err.Error()}, nil
		}

		if stock < quantity {
			return &pb.ValidateResponse{
				Valid:   false,
				Message: "Insufficient stock",
			}, nil
		}
	}

	if !found {
		return &pb.ValidateResponse{
			Valid:   false,
			Message: "Bundle not found",
		}, nil
	}

	return &pb.ValidateResponse{
		Valid:   true,
		Message: "Stock available",
	}, nil
}

// -----------------------------
// DEDUCT BUNDLE
// -----------------------------
func (s *server) DeductBundle(ctx context.Context, req *pb.BundleRequest) (*pb.DeductResponse, error) {

	tx, err := db.DB.Begin()
	if err != nil {
		return &pb.DeductResponse{Success: false, Message: err.Error()}, nil
	}

	query := `
	SELECT p.id, p.stock, bp.quantity
	FROM bundle_products bp
	JOIN products p ON p.id = bp.product_id
	WHERE bp.bundle_id = $1;
	`

	rows, err := tx.Query(query, req.BundleId)
	if err != nil {
		tx.Rollback()
		return &pb.DeductResponse{Success: false, Message: err.Error()}, nil
	}
	defer rows.Close()

	type Product struct {
		ID       string
		Stock    int
		Quantity int
	}

	var products []Product

	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Stock, &p.Quantity); err != nil {
			tx.Rollback()
			return &pb.DeductResponse{Success: false, Message: err.Error()}, nil
		}
		products = append(products, p)
	}

	if len(products) == 0 {
		tx.Rollback()
		return &pb.DeductResponse{Success: false, Message: "Bundle not found"}, nil
	}

	for _, p := range products {
		if p.Stock < p.Quantity {
			tx.Rollback()
			return &pb.DeductResponse{Success: false, Message: "Insufficient stock"}, nil
		}
	}

	for _, p := range products {
		_, err := tx.Exec(
			"UPDATE products SET stock = stock - $1 WHERE id = $2",
			p.Quantity, p.ID,
		)
		if err != nil {
			tx.Rollback()
			return &pb.DeductResponse{Success: false, Message: err.Error()}, nil
		}
	}

	if err := tx.Commit(); err != nil {
		return &pb.DeductResponse{Success: false, Message: err.Error()}, nil
	}

	return &pb.DeductResponse{
		Success: true,
		Message: "Stock deducted",
	}, nil
}

// -----------------------------
// VALIDATE STOCK
// -----------------------------
func (s *server) ValidateStock(ctx context.Context, req *pb.StockRequest) (*pb.StockValidateResponse, error) {
	var currentStock int

	query := `SELECT stock FROM products WHERE id = $1`
	err := db.DB.QueryRow(query, req.ProductId).Scan(&currentStock)

	if err != nil {
		return &pb.StockValidateResponse{
			Available:    false,
			CurrentStock: 0,
			Message:      "Product not found",
		}, nil
	}

	available := currentStock >= int(req.Quantity)
	message := "Stock available"
	if !available {
		message = "Insufficient stock"
	}

	return &pb.StockValidateResponse{
		Available:    available,
		CurrentStock: int32(currentStock),
		Message:      message,
	}, nil
}

// -----------------------------
// REDUCE STOCK
// -----------------------------
func (s *server) ReduceStock(ctx context.Context, req *pb.StockReductionRequest) (*pb.StockReductionResponse, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return &pb.StockReductionResponse{Success: false, Message: err.Error()}, nil
	}

	// First, validate all items have sufficient stock
	for _, item := range req.Items {
		var currentStock int
		err := tx.QueryRow("SELECT stock FROM products WHERE id = $1", item.ProductId).Scan(&currentStock)

		if err != nil {
			tx.Rollback()
			return &pb.StockReductionResponse{
				Success: false,
				Message: "Product " + item.ProductId + " not found",
			}, nil
		}

		if currentStock < int(item.Quantity) {
			tx.Rollback()
			return &pb.StockReductionResponse{
				Success: false,
				Message: "Insufficient stock for product " + item.ProductId,
			}, nil
		}
	}

	// If all validations pass, reduce the stock
	for _, item := range req.Items {
		_, err := tx.Exec(
			"UPDATE products SET stock = stock - $1 WHERE id = $2",
			item.Quantity, item.ProductId,
		)
		if err != nil {
			tx.Rollback()
			return &pb.StockReductionResponse{Success: false, Message: err.Error()}, nil
		}
	}

	if err := tx.Commit(); err != nil {
		return &pb.StockReductionResponse{Success: false, Message: err.Error()}, nil
	}

	return &pb.StockReductionResponse{
		Success: true,
		Message: "Stock reduced successfully",
	}, nil
}

// -----------------------------
// START gRPC SERVER
// -----------------------------
func StartGRPCServer() {

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterProductServiceServer(s, &server{})

	log.Println("gRPC server running on :50053")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
