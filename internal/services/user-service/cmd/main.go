package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/toji-dev/go-shop/internal/pkg/middleware"
	"github.com/toji-dev/go-shop/internal/pkg/tracing"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/container"
	user_grpc "github.com/toji-dev/go-shop/internal/services/user-service/internal/grpc"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/router"
	user_v1 "github.com/toji-dev/go-shop/proto/gen/go/user/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Debug configuration (chỉ trong development)
	if cfg.App.IsDevelopment() {
		cfg.Debug()
	}

	go func() {
		log.Println("Starting pprof server on :6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("Pprof server failed to start: %v", err)
		}
	}()

	jaegerAgentHost := "jaeger:4317"
	tp, err := tracing.InitTracerProvider(cfg.App.Name, jaegerAgentHost)
	if err != nil {
		log.Fatalf("failed to initialize tracer provider: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Initialize service container
	serviceContainer, err := container.NewServiceContainer(cfg)
	if err != nil {
		log.Fatal("Failed to initialize services:", err)
	}
	defer serviceContainer.Close()

	// Setup routes with service container
	router := router.SetupRoutes(serviceContainer)

	// Configure HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", serviceContainer.GetConfig().Server.Port),
		Handler:      router,
		ReadTimeout:  serviceContainer.GetConfig().Server.ReadTimeout,
		WriteTimeout: serviceContainer.GetConfig().Server.WriteTimeout,
		IdleTimeout:  serviceContainer.GetConfig().Server.IdleTimeout,
	}

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	go startMetricsServer()

	go runGrpcServer(serviceContainer.GetConfig(), &serviceContainer)

	// Start server in a goroutine
	go func() {
		log.Printf("%s starting on %s", serviceContainer.GetConfig().App.Name, serviceContainer.GetConfig().Server.GetServerAddress())
		log.Printf("Environment: %s", serviceContainer.GetConfig().App.Environment)
		log.Printf("Debug Mode: %v", serviceContainer.GetConfig().App.Debug)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
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

func runGrpcServer(cfg *config.Config, serviceContainer *container.ServiceContainer) {
	address := cfg.GRPC.Host + ":" + cfg.GRPC.Port
	log.Printf("Starting gRPC server on %s", address)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen for grpc on port %s: %v", cfg.GRPC.Host, err)
	}
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.PprofGRPCInterceptor(),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	grpcServer := user_grpc.NewUserGRPCServer(serviceContainer.GetAddressRepo())
	user_v1.RegisterUserServiceServer(s, grpcServer)
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
