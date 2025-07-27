package tracing

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitTracerProvider(jaegerURL, serviceName string) (*sdktrace.TracerProvider, func()) {
	ctx := context.Background()
	res, _ := resource.New(ctx, resource.WithAttributes(semconv.ServiceName(serviceName)))
	conn, _ := grpc.NewClient(jaegerURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	traceExporter, _ := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tracerProvider, func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			log.Printf("failed to shutdown tracer provider: %v", err)
		}
	}
}
