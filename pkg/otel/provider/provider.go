// Package provider は
package provider

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"golang.org/x/xerrors"
)

// Init はOTLP exporter を初期化し、対応するトレースプロバイダーとメトリックプロバイダーを設定する。
// https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/7cc24ffa9d815c71adf3802f58940a994e9ce9ec/examples/demo/server/main.go#L36
func Init(ctx context.Context, serviceNameKey string) (func(), error) {

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(serviceNameKey),
		),
	)
	if err != nil {
		return nil, xerrors.Errorf("failed to create resource. %w", err)
	}

	// TODO: この部分はパラメータでもらうようにする
	otelAgentAddr, ok := os.LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if !ok {
		otelAgentAddr = "0.0.0.0:4317"
	}

	metricExp, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(otelAgentAddr))
	if err != nil {
		return nil, xerrors.Errorf("failed to create the collector metric exporter. %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				metricExp,
				sdkmetric.WithInterval(60*time.Second), // TODO: この部分はパラメータでもらうようにする
			),
		),
	)
	otel.SetMeterProvider(meterProvider)

	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otelAgentAddr))

	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		return nil, xerrors.Errorf("failed to create the collector trace exporter. %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	return func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		if err := traceExp.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}

		// pushes any last exports to the receiver
		if err := meterProvider.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}, nil
}
