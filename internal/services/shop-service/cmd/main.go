package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"net/http"
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/pkg/middleware"
	"github.com/toji-dev/go-shop/internal/pkg/tracing"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/config"
	promotion_api "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/promotion/api"
	createpromotion "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/promotion/commands/create_promotion"
	deletepromotion "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/promotion/commands/delete_promotion"
	updatepromotion "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/promotion/commands/update_promotion"
	getpromotionbyid "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/promotion/queries/get_promotion_by_id"
	getpromotions "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/promotion/queries/get_promotions"
	shop_api "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/api"
	createshop "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/commands/create_shop"
	deleteshop "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/commands/delete_shop"
	updateshop "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/commands/update_shop"
	getshop "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/queries/get_shop"
	getshops "github.com/toji-dev/go-shop/internal/services/shop-service/internal/features/shop/queries/get_shops"
	shop_grpc "github.com/toji-dev/go-shop/internal/services/shop-service/internal/grpc"
	promotion_repo "github.com/toji-dev/go-shop/internal/services/shop-service/internal/repository/promotion"
	shop_repo "github.com/toji-dev/go-shop/internal/services/shop-service/internal/repository/shop"
	shop_v1 "github.com/toji-dev/go-shop/proto/gen/go/shop/v1"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func main() {

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

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

	go func() {
		log.Println("Starting pprof server on :6061")
		if err := http.ListenAndServe("localhost:6061", nil); err != nil {
			log.Printf("Pprof server failed to start: %v", err)
		}
	}()

	// Initialize database
	dbConfig := &postgresql_infra.DatabaseConfig{
		Host:         cfg.Database.Host,
		Port:         cfg.Database.Port,
		User:         cfg.Database.User,
		Password:     cfg.Database.Password,
		Name:         cfg.Database.DBName,
		SSLMode:      cfg.Database.SSLMode,
		MaxOpenConns: cfg.Database.MaxOpenConns,
		MaxIdleConns: cfg.Database.MaxIdleConns,
		MaxLifetime:  cfg.Database.MaxLifetime,
	}

	db, err := postgresql_infra.NewPostgreSQLService(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	shopRepo := shop_repo.NewPostgresShopRepository(db)
	promoRepo := promotion_repo.NewPostgresPromotionRepository(db)

	// Set Gin mode based on environment
	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	r := gin.Default()

	p := ginprometheus.NewPrometheus("go")
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.FullPath()
		if url == "" {
			url = "unknown"
		}
		return url
	}
	p.Use(r)

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize feature handlers
	// Create shop
	createShopHandler := createshop.NewHandler(shopRepo)
	createShopAPIHandler := createshop.NewAPIHandler(createShopHandler)

	// Get shop by ID
	getShopHandler := getshop.NewHandler(shopRepo)
	getShopAPIHandler := getshop.NewAPIHandler(getShopHandler)

	// Get shops by owner
	getShopsHandler := getshops.NewQueryHandler(shopRepo)
	getShopsAPIHandler := getshops.NewAPIHandler(getShopsHandler)

	// Update shop
	updateShopHandler := updateshop.NewCommandHandler(shopRepo)
	updateShopAPIHandler := updateshop.NewAPIHandler(updateShopHandler)

	// Delete shop
	deleteShopHandler := deleteshop.NewCommandHandler(shopRepo)
	deleteShopAPIHandler := deleteshop.NewAPIHandler(deleteShopHandler)

	// Create promotion
	createPromotionHandler := createpromotion.NewHandler(promoRepo)
	createPromotionApiHandler := createpromotion.NewAPIHandler(createPromotionHandler)

	// Get promotions
	getPromotionsHandler := getpromotions.NewHandler(promoRepo)
	getPromotionsAPIHandler := getpromotions.NewAPIHandler(getPromotionsHandler)

	// Get promotion by ID
	getPromotionByIDHandler := getpromotionbyid.NewHandler(promoRepo)
	getPromotionByIDAPIHandler := getpromotionbyid.NewAPIHandler(getPromotionByIDHandler)

	// Update promotion
	updatePromotionHandler := updatepromotion.NewHandler(promoRepo)
	updatePromotionAPIHandler := updatepromotion.NewAPIHandler(updatePromotionHandler)

	// Delete promotion
	deletePromotionHandler := deletepromotion.NewHandler(promoRepo)
	deletePromotionAPIHandler := deletepromotion.NewAPIHandler(deletePromotionHandler)

	// Register shop routes
	shop_api.RegisterShopRoutes(
		r,
		createShopAPIHandler,
		getShopAPIHandler,
		getShopsAPIHandler,
		updateShopAPIHandler,
		deleteShopAPIHandler,
	)

	promotion_api.RegisterPromotionRoutes(
		r,
		createPromotionApiHandler,
		getPromotionsAPIHandler,
		getPromotionByIDAPIHandler,
		updatePromotionAPIHandler,
		deletePromotionAPIHandler,
	)

	go startMetricsServer()

	go runGrpcServer(cfg, shopRepo, promoRepo)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Starting shop service on %s", cfg.GetServerAddress())
		if err := r.Run(cfg.GetServerAddress()); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down shop service...")
}

func runGrpcServer(config *config.Config, shopRepo shop_repo.ShopRepository, promotionRepo promotion_repo.PromotionRepository) {
	address := config.GRPC.Host + ":" + config.GRPC.Port
	log.Printf("Starting gRPC server on %s", address)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen for grpc on port %s: %v", config.GRPC.Host, err)
	}
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.PprofGRPCInterceptor(),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	grpcServer := shop_grpc.NewShopGRPCServer(shopRepo, promotionRepo)
	shop_v1.RegisterShopServiceServer(s, grpcServer)
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}

func startMetricsServer() {
	metricsRouter := gin.New()
	metricsRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))
	log.Println("Starting metrics server on 0.0.0.0:9100")
	if err := metricsRouter.Run("0.0.0.0:9100"); err != nil {
		log.Fatalf("Failed to start metrics server: %v", err)
	}
}
