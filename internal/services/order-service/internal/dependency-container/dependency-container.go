package dependency_container

import (
	"fmt"
	"log"

	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/grpc/adapter"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/handler"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/repository"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/usecase"
)

type DependencyContainer struct {
	config       *config.Config
	postgreSQL   *postgresql_infra.PostgreSQLService
	orderRepo    repository.OrderRepository
	orderUsecase usecase.OrderUsecase
	orderHandler handler.OrderHandler

	shopServiceAdapter    adapter.ShopServiceAdapter
	productServiceAdapter adapter.ProductServiceAdapter
	userAdapter           adapter.UserServiceAdapter
}

func NewDependencyContainer(cfg *config.Config) *DependencyContainer {
	container := &DependencyContainer{
		config: cfg,
	}

	if err := container.initPostgreSQL(); err != nil {
		log.Fatalf("failed to initialize PostgreSQL: %v", err)
	}

	container.initRepositories()

	container.initShopServiceAdapter()

	container.initProductServiceAdapter()

	container.initUserServiceAdapter()

	container.initUseCases()

	container.initOrderHandler()

	return container
}

func (sc *DependencyContainer) initPostgreSQL() error {
	pgConfig := &postgresql_infra.DatabaseConfig{
		Host:         sc.config.Database.Host,
		Port:         sc.config.Database.Port,
		User:         sc.config.Database.User,
		Password:     sc.config.Database.Password,
		Name:         sc.config.Database.DBName,
		SSLMode:      sc.config.Database.SSLMode,
		MaxOpenConns: sc.config.Database.MaxOpenConns,

		MaxIdleConns: sc.config.Database.MaxIdleConns,
		MaxLifetime:  sc.config.Database.MaxLifetime,
	}

	pgService, err := postgresql_infra.NewPostgreSQLService(pgConfig)
	if err != nil {
		return fmt.Errorf("failed to create PostgreSQL service: %w", err)
	}

	sc.postgreSQL = pgService
	log.Println("PostgreSQL service initialized")
	return nil
}

func (sc *DependencyContainer) initRepositories() {
	sc.orderRepo = repository.NewOrderRepository(sc.postgreSQL)
	log.Println("Order repository initialized")
}

func (sc *DependencyContainer) initUseCases() {
	sc.orderUsecase = usecase.NewOrderUsecase(
		sc.orderRepo,
		sc.shopServiceAdapter,
		sc.productServiceAdapter,
		sc.userAdapter,
	)
	log.Println("Order use case initialized")
}

func (sc *DependencyContainer) initOrderHandler() {
	sc.orderHandler = handler.NewOrderHandler(sc.orderUsecase)
	log.Println("Order handler initialized")
}

func (sc *DependencyContainer) initShopServiceAdapter() error {
	shopServiceAddr := fmt.Sprintf("%s:%s", sc.config.ShopServiceAdapter.Host, sc.config.ShopServiceAdapter.Port)
	if shopServiceAddr == "" {
		return fmt.Errorf("shop service address is not configured")
	}

	adapter, err := adapter.NewGrpcShopAdapter(shopServiceAddr)
	if err != nil {
		return fmt.Errorf("failed to create shop service adapter: %w", err)
	}

	sc.shopServiceAdapter = adapter
	log.Println("Shop service adapter initialized")
	return nil
}

func (sc *DependencyContainer) initProductServiceAdapter() error {
	productServiceAddr := fmt.Sprintf("%s:%s", sc.config.ProductServiceAdapter.Host, sc.config.ProductServiceAdapter.Port)
	if productServiceAddr == "" {
		return fmt.Errorf("product service address is not configured")
	}

	adapter, err := adapter.NewGrpcProductAdapter(productServiceAddr)
	if err != nil {
		return fmt.Errorf("failed to create product service adapter: %w", err)
	}

	sc.productServiceAdapter = adapter
	log.Println("Product service adapter initialized")
	return nil
}

func (sc *DependencyContainer) initUserServiceAdapter() error {
	userServiceAddr := fmt.Sprintf("%s:%s", sc.config.UserServiceAdapter.Host, sc.config.UserServiceAdapter.Port)
	if userServiceAddr == "" {
		return fmt.Errorf("user service address is not configured")
	}

	adapter, err := adapter.NewGrpcUserAdapter(userServiceAddr)
	if err != nil {
		return fmt.Errorf("failed to create user service adapter: %w", err)
	}

	sc.userAdapter = adapter
	log.Println("User service adapter initialized")
	return nil
}

func (sc *DependencyContainer) GetOrderHandler() handler.OrderHandler {
	return sc.orderHandler
}

func (sc *DependencyContainer) GetConfig() *config.Config {
	return sc.config
}

func (sc *DependencyContainer) GetOrderRepository() repository.OrderRepository {
	return sc.orderRepo
}

func (sc *DependencyContainer) GetProductServiceAdapter() adapter.ProductServiceAdapter {
	return sc.productServiceAdapter
}
