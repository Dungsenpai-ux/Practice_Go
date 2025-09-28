# Dịch vụ Practice Go

Một dịch vụ HTTP API nhỏ viết bằng Go theo mô hình nhiều tầng (layered) minh hoạ các ý tưởng:

- Kiến trúc phân tầng (controller / service / repository / model / middleware)
- Lưu trữ PostgreSQL + migrations
- Cache với Memcached (cache dương + cache âm)
- Tracing OpenTelemetry gửi trực tiếp Tempo (xem trên Grafana)
- Metrics đơn giản bằng expvar (đã gỡ Prometheus) + logging cơ bản

## Cấu trúc thư mục

```
config/              # Nạp cấu hình & biến môi trường
controller/          # HTTP handlers (Gin + tương thích net/http)
model/               # Domain models (Movie, Health)
service/             # Business logic, tích hợp ngoài
service/db/          # Repository cho Postgres (pgx)
service/external/    # Dịch vụ bên ngoài (memcache wrapper...)
service/middleware/  # Tracing, logging, expvar helpers
service/observability/tempo.yaml/  # Cấu hình Tempo (config.yaml bên trong)
docs/                # Tài liệu dự án
```

## Tính năng chính

| Thành phần | Triển khai |
|------------|------------|
| HTTP Framework | Gin (otelgin middleware) |
| Tracing | OpenTelemetry SDK → Tempo (OTLP HTTP hoặc gRPC) |
| Metrics | expvar (cung cấp ở /metrics & /debug/vars) |
| Database | PostgreSQL (pgx pool) |
| Cache | Memcached (read-through + negative caching) |
| Migrations | File SQL trong `service/middleware/migrations` (chạy bằng golang-migrate bên ngoài) |
| Config | Biến môi trường parse trong `config/config.go` |

## Khởi động nhanh

1. Sao chép `.env.example` thành `.env` (nếu có) và chỉnh sửa.
2. Khởi động stack tracing (Tempo + Grafana):
   ```bash
   docker compose -f docker-compose.tracing.yml up -d
   ```
3. Chạy ứng dụng:
   ```bash
   go run ./...
   ```
4. Gọi endpoint để sinh trace:
   ```bash
   curl http://localhost:8080/healthz
   ```
5. Mở Grafana → Explore → chọn datasource Tempo để xem trace.

## Tóm tắt API

| Method | Đường dẫn | Mô tả |
|--------|-----------|-------|
| GET | /healthz | Thông tin tình trạng & version |
| POST | /movies | Tạo phim (JSON body) |
| GET | /movies/{id} | Lấy phim theo ID (có cache) |
| GET | /movies/search?q=TERM&year=YYYY | Tìm kiếm phim |

Xem chi tiết request/response trong `api.md`.

## Observability
- Trace tự động có đầy đủ resource attributes.
- Có thể mở rộng fallback exporter ra stdout nếu Tempo down.
- Metrics (expvar) hiển thị tại `/metrics` và `/debug/vars`.

## Hướng phát triển thêm
- Thêm retry & circuit breaker.
- Thêm integration test với Postgres test container.
- Thêm structured logging (zap hoặc zerolog).
- Docker hoá chính ứng dụng (service app) nếu cần triển khai.

---
Các tài liệu liên quan:
- `setup.md` (Hướng dẫn cài đặt)
- `configuration.md` (Cấu hình)
- `observability.md` (Giám sát)
- `api.md` (API chi tiết)
- `data-model.md` (Mô hình dữ liệu)
- `architecture.md` (Kiến trúc)
- `contributing.md` (Đóng góp)
