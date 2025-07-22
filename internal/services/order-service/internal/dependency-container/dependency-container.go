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

	shopServiceAdapter adapter.ShopServiceAdapter
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
	sc.orderUsecase = usecase.NewOrderUsecase(sc.orderRepo, sc.shopServiceAdapter)
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

func (sc *DependencyContainer) GetOrderHandler() handler.OrderHandler {
	return sc.orderHandler
}

func (sc *DependencyContainer) GetConfig() *config.Config {
	return sc.config
}
