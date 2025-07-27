package dependency_container

import (
	"fmt"
	"log"

	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/config"
	grpc_adapter "github.com/toji-dev/go-shop/internal/services/payment-service/internal/grpc/adapter"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/handler"
	paymentprovider "github.com/toji-dev/go-shop/internal/services/payment-service/internal/payment_provider"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/repository"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/usecase"
)

type DependencyContainer struct {
	config         *config.Config
	postgreSQL     *postgresql_infra.PostgreSQLService
	paymentRepo    repository.PaymentRepository
	paymentUseCase usecase.PaymentUseCase
	paymentHandler handler.PaymentHandler

	paymentMethodFactory *paymentprovider.PaymentProviderFactory

	orderServiceAdapter grpc_adapter.OrderServiceAdapter
}

func NewDependencyContainer(cfg *config.Config) *DependencyContainer {
	container := &DependencyContainer{
		config: cfg,
	}

	if err := container.initPostgreSQL(); err != nil {
		log.Fatalf("failed to initialize PostgreSQL: %v", err)
	}

	container.initRepositories()

	container.initPaymentProviders()

	container.initOrderServiceAdapter()

	container.initUseCases()

	container.initPaymentHandler()

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
	sc.paymentRepo = repository.NewPaymentRepository(sc.postgreSQL)
	log.Println("Payment repository initialized")
}

func (sc *DependencyContainer) initPaymentProviders() {
	sc.paymentMethodFactory = paymentprovider.NewPaymentProviderFactory()
	log.Println("Payment provider factory initialized")
}

func (sc *DependencyContainer) initOrderServiceAdapter() {
	var err error
	address := fmt.Sprintf("%s:%d", sc.config.OrderGrpcConfig.OrderServiceHost, sc.config.OrderGrpcConfig.OrderServicePort)
	sc.orderServiceAdapter, err = grpc_adapter.NewGrpcOrderAdapter(address)
	if err != nil {
		log.Fatalf("failed to initialize order service adapter: %v", err)
	}
	log.Println("Order service adapter initialized")
}

func (sc *DependencyContainer) initUseCases() {
	sc.paymentUseCase = usecase.NewPaymentUsecase(
		sc.paymentRepo,
		sc.paymentMethodFactory,
		sc.orderServiceAdapter,
	)
	log.Println("Payment use case initialized")
}

func (sc *DependencyContainer) initPaymentHandler() {
	sc.paymentHandler = handler.NewPaymentHandler(sc.paymentUseCase)
	log.Println("Payment handler initialized")
}

func (sc *DependencyContainer) GetPaymentHandler() handler.PaymentHandler {
	return sc.paymentHandler
}

func (sc *DependencyContainer) GetConfig() *config.Config {
	return sc.config
}

func (sc *DependencyContainer) GetPaymentRepository() repository.PaymentRepository {
	return sc.paymentRepo
}
