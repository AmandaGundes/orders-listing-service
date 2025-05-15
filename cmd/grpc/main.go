package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	pb "github.com/agundes/Projects/go/orders-listing-service/proto"
)

type server struct {
	pb.UnimplementedOrderServiceServer
	db *sql.DB
}

func (s *server) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	rows, err := s.db.Query("SELECT id, customer_name, created_at FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resp := &pb.ListOrdersResponse{}
	for rows.Next() {
		var o pb.Order
		if err := rows.Scan(&o.Id, &o.CustomerName, &o.CreatedAt); err != nil {
			return nil, err
		}
		resp.Orders = append(resp.Orders, &o)
	}
	return resp, nil
}

func main() {
	dsn := os.Getenv("DB_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, &server{db: db})
	log.Println("gRPC service listening on :50051")
	log.Fatal(s.Serve(lis))
}
