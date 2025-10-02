package middleware

import (
	"context"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Dungsenpai-ux/Practice_Go/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// InitTracer configures global tracer provider (OTLP/HTTP) compatible with Tempo.
func InitTracer(cfg *config.Config) func(context.Context) error {
	ctx := context.Background()

	endpoint := sanitizeEndpoint(cfg.OtelEndpoint)
	protocol := strings.ToLower(os.Getenv("OTEL_EXPORTER_OTLP_TRACES_PROTOCOL"))
	var exporter sdktrace.SpanExporter
	var err error
	switch protocol {
	case "grpc":
		hostPort := endpoint
		if u, perr := url.Parse(endpoint); perr == nil && u.Host != "" { // strip scheme
			hostPort = u.Host
		}
		exporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(hostPort),
			otlptracegrpc.WithInsecure(),
		)
	default: // http
		exporter, err = otlptracehttp.New(ctx,
			otlptracehttp.WithEndpointURL(endpoint),
			otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
		)
		protocol = "http"
	}
	if err != nil {
		log.Printf("[otel] primary exporter failed (%s): %v -- falling back to stdout", protocol, err)
		if se, serr := stdouttrace.New(stdouttrace.WithPrettyPrint()); serr == nil {
			exporter = se
			protocol = "stdout"
		} else {
			log.Printf("[otel] stdout exporter failed: %v (tracing disabled)", serr)
			return func(context.Context) error { return nil }
		}
	}

	res, rerr := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.OtelService),
			semconv.ServiceVersion(cfg.Version),
			attribute.String("deployment.environment", cfg.OtelEnv),
		),
	)
	if rerr != nil {
		log.Printf("[otel] resource init failed: %v", rerr)
	}

	sampler := parseSampler(cfg.OtelSampler)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(5*time.Second)),
	)
	otel.SetTracerProvider(tp)
	log.Printf("[otel] tracing enabled protocol=%s endpoint=%s service=%s sampler=%s", protocol, endpoint, cfg.OtelService, cfg.OtelSampler)
	return tp.Shutdown
}

func sanitizeEndpoint(ep string) string {
	if strings.HasSuffix(ep, "/v1/traces") {
		return strings.TrimSuffix(ep, "/v1/traces")
	}
	return ep
}

func parseSampler(spec string) sdktrace.Sampler {
	s := strings.ToLower(strings.TrimSpace(spec))
	switch {
	case s == "always_off":
		return sdktrace.NeverSample()
	case s == "always_on":
		return sdktrace.AlwaysSample()
	case s == "parentbased_always_on":
		return sdktrace.ParentBased(sdktrace.AlwaysSample())
	case strings.HasPrefix(s, "ratio:"):
		parts := strings.SplitN(s, ":", 2)
		if len(parts) == 2 {
			if v, err := strconv.ParseFloat(parts[1], 64); err == nil && v >= 0 && v <= 1 {
				return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(v))
			}
		}
		return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(1.0))
	default:
		return sdktrace.ParentBased(sdktrace.AlwaysSample())
	}
}

// Tracer helper returns a named tracer for custom spans.
func Tracer() trace.Tracer { return otel.Tracer("app") }

// ForceFlush flushes spans (best-effort) – useful on short-lived commands/tests.
func ForceFlush(ctx context.Context) {
	if tp, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider); ok {
		_ = tp.ForceFlush(ctx)
	}
}
