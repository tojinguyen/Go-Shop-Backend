package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

func OtelTracingMiddleware(serviceName string) gin.HandlerFunc {
	tracer := otel.Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(c *gin.Context) {
		// Trích xuất context từ header của request đến
		ctx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		// Tạo một span mới
		spanName := c.FullPath()
		if spanName == "" {
			spanName = c.Request.URL.Path
		}

		ctx, span := tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		// Thêm các thuộc tính tiêu chuẩn vào span
		span.SetAttributes(
			semconv.HTTPMethodKey.String(c.Request.Method),
			semconv.HTTPURLKey.String(c.Request.URL.String()),
			semconv.NetHostNameKey.String(c.Request.Host),
			semconv.HTTPRequestContentLengthKey.Int64(c.Request.ContentLength),
			attribute.String("http.client_ip", c.ClientIP()),
		)

		// Gắn context mới vào Gin context để các handler sau có thể sử dụng
		c.Request = c.Request.WithContext(ctx)

		// Đi tiếp
		c.Next()

		// Ghi lại status code sau khi request hoàn thành
		span.SetAttributes(semconv.HTTPStatusCodeKey.Int(c.Writer.Status()))
		if c.Writer.Status() >= 500 {
			span.SetStatus(codes.Error, "Server Error")
		}
	}
}
