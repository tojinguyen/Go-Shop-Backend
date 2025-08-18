package middleware

import (
	"context"
	"runtime/pprof"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// PprofGRPCInterceptor là một UnaryServerInterceptor để gắn nhãn trace_id vào pprof.
func PprofGRPCInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Lấy span từ context mà OpenTelemetry đã truyền vào
		span := trace.SpanFromContext(ctx)

		// Kiểm tra xem span có hợp lệ và đang được ghi lại không
		if span.IsRecording() {
			traceID := span.SpanContext().TraceID().String()

			// Tạo một context mới với nhãn pprof
			labels := pprof.Labels("trace_id", traceID)
			ctxWithLabels := pprof.WithLabels(ctx, labels)

			// Thực thi handler với context đã được gắn nhãn
			pprof.Do(ctxWithLabels, labels, func(ctx context.Context) {
				// Handler sẽ được chạy trong scope của pprof.Do
			})

			// Gọi handler gốc với context đã được gắn nhãn
			return handler(ctxWithLabels, req)
		}

		// Nếu không có trace, chỉ cần gọi handler gốc
		return handler(ctx, req)
	}
}
