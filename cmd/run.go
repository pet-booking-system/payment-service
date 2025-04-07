package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"payment-service/config"
	"payment-service/internal/logger"
	"payment-service/internal/repository"
	"payment-service/internal/server"
	"payment-service/internal/service"

	paymentpb "github.com/pet-booking-system/proto-definitions/payment"

	"google.golang.org/grpc"
)

func Run() {
	logger.Init()

	cfg := config.Load()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	repo := repository.NewPaymentRepository(db)
	svc := service.NewPaymentService(repo)
	handler := server.NewPaymentHandler(svc)

	grpcServer := grpc.NewServer()
	paymentpb.RegisterPaymentServiceServer(grpcServer, handler)

	addr := fmt.Sprintf(":%s", cfg.GRPCPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", cfg.GRPCPort, err)
	}

	logger.Info("Payment service started on port ", cfg.GRPCPort)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
