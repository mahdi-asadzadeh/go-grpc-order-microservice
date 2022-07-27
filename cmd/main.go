package main

import (
	"fmt"
	"github.com/mahdi-asadzadeh/go-grpc-order-microservice/pkg/client"
	"github.com/mahdi-asadzadeh/go-grpc-order-microservice/pkg/config"
	"github.com/mahdi-asadzadeh/go-grpc-order-microservice/pkg/db"
	"github.com/mahdi-asadzadeh/go-grpc-order-microservice/pkg/pb"
	"github.com/mahdi-asadzadeh/go-grpc-order-microservice/pkg/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("Starting order server ...")
	// Load config
	config.LoadSettings(true)

	// Database config
	DB_URL := os.Getenv("DB_URL")
	h := db.Init(DB_URL)

	lis, err := net.Listen("tcp", os.Getenv("SERVER_IP"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Register product service
	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, &service.Server{H: h, ProductClient: client.InitServiceClient()})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
