// Package observability - tracing.go: OpenTelemetry SDK 초기화.
// 환경변수 기반 설정 (OTEL_EXPORTER_OTLP_ENDPOINT, OTEL_SERVICE_NAME 등) 을
// SDK 가 자동 인식. main 에서 InitTracing 1회 호출 후 shutdown 을 defer.
package observability

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

// InitTracing 은 TracerProvider 와 OTLP gRPC exporter 를 설정하고
// global tracer 로 등록한다. 반환된 shutdown 을 프로세스 종료 시 호출하여
// 버퍼 flush 가 일어나도록 한다.
func InitTracing(ctx context.Context, serviceName string) (func(context.Context) error, error) {
	exp, err := otlptrace.New(ctx, otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
	))
	if err != nil {
		return nil, err
	}

	res, err := sdkresource.New(ctx,
		sdkresource.WithFromEnv(),
		sdkresource.WithProcess(),
		sdkresource.WithTelemetrySDK(),
		sdkresource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp, sdktrace.WithBatchTimeout(5*time.Second)),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown, nil
}
