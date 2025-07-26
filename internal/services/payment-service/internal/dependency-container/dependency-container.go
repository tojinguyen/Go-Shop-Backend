package dependency_container

import (
	"fmt"
	"log"

	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/config"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/handler"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/repository"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/usecase"
)

type DependencyContainer struct {
	config       *config.Config
	postgreSQL   *postgresql_infra.PostgreSQLService
	orderRepo    repository.PaymentRepository
	orderUsecase usecase.PaymentUseCase
	orderHandler handler.PaymentHandler
}

func NewDependencyContainer(cfg *config.Config) *DependencyContainer {
	container := &DependencyContainer{
		config: cfg,
	}

	if err := container.initPostgreSQL(); err != nil {
		log.Fatalf("failed to initialize PostgreSQL: %v", err)
	}

	container.initRepositories()

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
	sc.orderRepo = repository.NewPaymentRepository(sc.postgreSQL)
	log.Println("Order repository initialized")
}

func (sc *DependencyContainer) initUseCases() {
	sc.orderUsecase = usecase.NewPaymentUsecase(
		sc.orderRepo,
	)
	log.Println("Order use case initialized")
}

func (sc *DependencyContainer) initOrderHandler() {
	sc.orderHandler = handler.NewOrderHandler(sc.orderUsecase)
	log.Println("Order handler initialized")
}

func (sc *DependencyContainer) GetOrderHandler() handler.PaymentHandler {
	return sc.orderHandler
}

func (sc *DependencyContainer) GetConfig() *config.Config {
	return sc.config
}

func (sc *DependencyContainer) GetPaymentRepository() repository.PaymentRepository {
	return sc.orderRepo
}
