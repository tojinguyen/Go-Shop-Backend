package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	product_grpc "github.com/toji-dev/go-shop/internal/services/product-service/internal/grpc"
	product_v1 "github.com/toji-dev/go-shop/proto/gen/go/product/v1"
	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/config"
	dependency_container "github.com/toji-dev/go-shop/internal/services/product-service/internal/dependency-container"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/repository"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/router"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	dependencyContainer, err := dependency_container.NewDependencyContainer(cfg)

	if err != nil {
		log.Fatalf("Failed to initialize dependency container: %v", err)
	}
	defer dependencyContainer.Close()

	router := router.SetupRoutes(r, dependencyContainer)

	server := &http.Server{
		Addr:         dependencyContainer.GetConfig().Server.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  dependencyContainer.GetConfig().Server.ReadTimeout,
		WriteTimeout: dependencyContainer.GetConfig().Server.WriteTimeout,
		IdleTimeout:  dependencyContainer.GetConfig().Server.IdleTimeout,
	}

	go startMetricsServer()

	go runGrpcServer(cfg, dependencyContainer.GetProductRepository())

	go func() {
		log.Printf("%s starting on %s", dependencyContainer.GetConfig().App.Name, dependencyContainer.GetConfig().Server.GetServerAddress())
		log.Printf("Environment: %s", dependencyContainer.GetConfig().App.Environment)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests a deadline of 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func startMetricsServer() {
	metricsRouter := gin.New()
	metricsRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))
	log.Println("Starting metrics server on 0.0.0.0:9100")
	if err := metricsRouter.Run("0.0.0.0:9100"); err != nil {
		log.Fatalf("Failed to start metrics server: %v", err)
	}
}

func runGrpcServer(cfg *config.Config, productRepo repository.ProductRepository) {
	address := cfg.GRPC.Host + ":" + cfg.GRPC.Port
	log.Printf("Starting gRPC server on %s", address)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen for grpc on port %s: %v", cfg.GRPC.Host, err)
	}
	s := grpc.NewServer()
	grpcServer := product_grpc.NewProductGRPCServer(productRepo)
	product_v1.RegisterProductServiceServer(s, grpcServer)
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
