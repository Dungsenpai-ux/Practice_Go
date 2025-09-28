# Khả năng quan sát (Observability)

Dịch vụ này hỗ trợ tracing (OpenTelemetry) và các bộ đếm đơn giản bằng expvar.

## Chuỗi xử lý Tracing
Application → OpenTelemetry SDK → OTLP HTTP/gRPC → Tempo → Grafana Explore.

### Chọn Exporter
Mã kiểm tra biến `OTEL_EXPORTER_OTLP_TRACES_PROTOCOL`:
- `grpc` → dùng OTLP gRPC exporter tới endpoint:4317
- khác (`http` hoặc trống) → dùng OTLP HTTP exporter tới endpoint:4318
Nếu khởi tạo thất bại, có thể bật fallback stdout exporter (xem mã middleware tracing).

### Thuộc tính Resource
Bao gồm:
- service.name
- service.version (nếu cấu hình)
- deployment.environment (APP_ENV)
- host.*, os.*, process.* (thuộc tính chuẩn)

### Sampler
Cấu hình qua `OTEL_TRACES_SAMPLER` hỗ trợ lấy mẫu theo tỉ lệ tùy chỉnh (định dạng `ratio:0.3`).

### Xem Traces
1. Mở Grafana → Explore
2. Chọn datasource Tempo (mặc định)
3. Chọn khoảng thời gian gần đây và chạy truy vấn
4. Lọc theo Service Name = OTEL_SERVICE_NAME của bạn

## Metric expvar
Endpoints:
- `/metrics` (JSON expvar)
- `/debug/vars` (nội dung tương tự)

Các bộ đếm chính (ví dụ):
- http_requests_total
- http_requests_in_flight
- http_request_duration_ms_sum (tổng hợp map)
- http_status_count
- cache_hits / cache_misses / cache_errors

Các số liệu này chưa ở định dạng Prometheus; sau này bạn có thể bọc lại bằng collector hoặc quay lại dùng Prometheus.

## Bổ sung tương quan Logs (Tương lai)
- Thêm logger có cấu trúc (zap) và chèn trace/span ID từ context.

## Nâng cấp tiềm năng
| Khu vực | Ý tưởng |
|--------|---------|
| Metrics | Đưa lại Prometheus hoặc OTEL metrics SDK |
| Tracing | Tinh chỉnh batch span processor / resource detectors |
| Logs | Thêm pipeline OTEL logs |
| Profiling | Tích hợp các endpoint pprof |
