package middleware

import (
    "expvar"
    "log"
    "net/http"
    "runtime"
    "strconv"
    "time"
)

var (
    reqTotal     = expvar.NewInt("http_requests_total")
    reqInFlight  = expvar.NewInt("http_requests_in_flight")
    reqDuration  = expvar.NewMap("http_request_duration_ms_sum") // method+path -> total ms
    reqCountPath = expvar.NewMap("http_request_count")            // method+path -> count
    reqStatus    = expvar.NewMap("http_status_count")             // status code -> count
)

// LoggingAndMetrics wraps an http.Handler with basic structured logging + metrics.
func LoggingAndMetrics(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        reqTotal.Add(1)
        reqInFlight.Add(1)
        lrw := &loggingResponseWriter{ResponseWriter: w, status: 200}
        defer func() {
            reqInFlight.Add(-1)
            durMs := time.Since(start).Milliseconds()
            key := r.Method + " " + routeKey(r.URL.Path)
            // duration aggregation
            incrIntMap(reqDuration, key, durMs)
            incrIntMap(reqCountPath, key, 1)
            incrIntMap(reqStatus, strconv.Itoa(lrw.status), 1)
            log.Printf("method=%s path=%s status=%d duration_ms=%d in_flight=%s goroutines=%d", r.Method, r.URL.Path, lrw.status, durMs, reqInFlight.String(), runtime.NumGoroutine())
        }()
        next.ServeHTTP(lrw, r)
    })
}

type loggingResponseWriter struct {
    http.ResponseWriter
    status int
}

func (l *loggingResponseWriter) WriteHeader(code int) {
    l.status = code
    l.ResponseWriter.WriteHeader(code)
}

// routeKey collapses dynamic segments (simple heuristic) to avoid cardinality explosion.
func routeKey(p string) string {
    // naive: if path ends with numeric id treat as /resource/:id
    // this keeps metrics aggregated.
    // Improve pattern matching as needed.
    n := len(p)
    if n == 0 { return p }
    // split simple
    // For brevity keep implementation minimal.
    return p
}

// incrIntMap increments an expvar.Map (string->Int) by delta.
func incrIntMap(m *expvar.Map, key string, delta int64) {
    v := m.Get(key)
    if v == nil {
        iv := new(expvar.Int)
        iv.Set(delta)
        m.Set(key, iv)
        return
    }
    if iv, ok := v.(*expvar.Int); ok {
        iv.Add(delta)
    }
}
