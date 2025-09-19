package config

import (
	"net/http"
	"os"
	"sync"
	"time"

	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	once sync.Once

	// Metrics
	HttpRequestsTotal   *prometheus.CounterVec
	HttpRequestDuration *prometheus.HistogramVec
	DBQueryDuration     prometheus.Histogram
	CacheHit            prometheus.Counter
	CacheMiss           prometheus.Counter
)

// Init registers metrics. Safe to call multiple times.
func Init() {
	once.Do(func() {
		HttpRequestsTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		)
		HttpRequestDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Histogram of latencies for HTTP requests",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		)
		DBQueryDuration = prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "Histogram of database query latencies",
				Buckets: prometheus.DefBuckets,
			},
		)
		CacheHit = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "cache_hit_total",
				Help: "Total cache hits",
			},
		)
		CacheMiss = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "cache_miss_total",
				Help: "Total cache misses",
			},
		)

		prometheus.MustRegister(HttpRequestsTotal, HttpRequestDuration, DBQueryDuration, CacheHit, CacheMiss)
	})
}

// SetupLogger configures zerolog for JSON output.
func SetupLogger() {
	// Zerolog is JSON by default
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

// statusRecorder wraps ResponseWriter to capture status code
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// InstrumentHandlerWithPath instruments an http.Handler with metrics and logging.
func InstrumentHandlerWithPath(path string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if HttpRequestsTotal == nil {
			Init()
		}
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: 200}
		next.ServeHTTP(rec, r)
		duration := time.Since(start).Seconds()
		statusLabel := strconv.Itoa(rec.status)
		HttpRequestsTotal.WithLabelValues(r.Method, path, statusLabel).Inc()
		HttpRequestDuration.WithLabelValues(r.Method, path, statusLabel).Observe(duration)
	})
}

// InstrumentHandler instruments without a specific path (uses r.URL.Path). Prefer WithPath.
func InstrumentHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if HttpRequestsTotal == nil {
			Init()
		}
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: 200}
		next.ServeHTTP(rec, r)
		duration := time.Since(start).Seconds()
		statusLabel := strconv.Itoa(rec.status)
		path := r.URL.Path
		HttpRequestsTotal.WithLabelValues(r.Method, path, statusLabel).Inc()
		HttpRequestDuration.WithLabelValues(r.Method, path, statusLabel).Observe(duration)
	})
}

// ObserveDBQueryDuration returns a stop function which records elapsed time.
func ObserveDBQueryDuration() func() {
	if DBQueryDuration == nil {
		Init()
	}
	start := time.Now()
	return func() { DBQueryDuration.Observe(time.Since(start).Seconds()) }
}
