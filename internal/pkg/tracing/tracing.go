package tracing

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitTracerProvider(serviceName, jaegerAgentHost string) (*sdktrace.TracerProvider, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Thiết lập kết nối gRPC đến Jaeger Agent/Collector
	conn, err := grpc.DialContext(ctx, jaegerAgentHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	// Tạo một OTLP trace exporter
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	// Xác định các thuộc tính của service (tên, phiên bản, môi trường)
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	// Tạo TracerProvider với BatchSpanProcessor để gửi span theo lô.
	// Điều này hiệu quả hơn việc gửi từng span một.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Đăng ký TracerProvider làm provider global
	otel.SetTracerProvider(tp)

	// Đăng ký TextMapPropagator để context có thể được truyền qua các request.
	// Đây là phần quan trọng để liên kết các trace giữa các service.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	log.Printf("Tracer provider for service '%s' initialized and registered.", serviceName)
	return tp, nil
}
