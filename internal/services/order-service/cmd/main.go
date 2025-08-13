package main

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/config"
	dependency_container "github.com/toji-dev/go-shop/internal/services/order-service/internal/dependency-container"
	grpc_server "github.com/toji-dev/go-shop/internal/services/order-service/internal/grpc/server"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/repository"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/router"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/worker"
	order_v1 "github.com/toji-dev/go-shop/proto/gen/go/order/v1"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		return
	}
	fmt.Println("Starting order-service...")

	// Initialize dependency container
	dependencyContainer := dependency_container.NewDependencyContainer(cfg)

	scheduler := worker.NewScheduler(dependencyContainer)
	scheduler.RegisterJobs()

	// Start Kafka consumer with inbox pattern
	kafkaConsumer := worker.NewKafkaConsumer(
		cfg,
		dependencyContainer.GetOrderUsecase(),
		dependencyContainer.GetInboxEventUsecase(),
	)

	// Start inbox processing worker
	inboxWorker := worker.NewInboxWorker(dependencyContainer.GetInboxEventUsecase())

	go startMetricsServer()

	go kafkaConsumer.StartAllKafkaConsumer()
	go inboxWorker.Start()

	// Chạy scheduler trong một goroutine riêng
	go scheduler.Start()

	go runGrpcServer(cfg, dependencyContainer.GetOrderRepository())

	// Initialize Gin router
	r := gin.Default()

	// Setup routes
	router.Init(r, dependencyContainer)

	// Start server
	if err := r.Run(cfg.Server.GetServerAddress()); err != nil {
		fmt.Printf("Failed to start server: %v", err)
	}
}

func runGrpcServer(cfg *config.Config, orderRepo repository.OrderRepository) {
	address := cfg.GRPC.ServiceHost + ":" + strconv.Itoa(cfg.GRPC.ServicePort)
	log.Printf("Starting gRPC server on %s", address)
	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("failed to listen for grpc on port %d: %v", cfg.GRPC.ServicePort, err)
	}

	s := grpc.NewServer()
	server := grpc_server.NewOrderGRPCServer(orderRepo)

	order_v1.RegisterOrderServiceServer(s, server)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}

	log.Printf("gRPC server listening at %v", lis.Addr())
}

func startMetricsServer() {
	metricsRouter := gin.New()
	metricsRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))
	log.Println("Starting metrics server on :9100")
	if err := metricsRouter.Run(":9100"); err != nil {
		log.Fatalf("Failed to start metrics server: %v", err)
	}
}
