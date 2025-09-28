# Cấu hình

Cấu hình ứng dụng được nạp từ các biến môi trường trong `config/config.go`.

## Biến lõi
| Tên | Mục đích | Ghi chú |
|-----|----------|---------|
| APP_ENV | Môi trường (local/dev/staging/prod) | Thêm vào resource attributes |
| OTEL_EXPORTER_OTLP_ENDPOINT | Endpoint OTLP gốc (không thêm /v1/traces) | Hỗ trợ HTTP hoặc gRPC (4318/4317) |
| OTEL_EXPORTER_OTLP_TRACES_PROTOCOL | http hoặc grpc | Tùy chọn; mặc định http nếu bỏ trống |
| OTEL_SERVICE_NAME | Tên logic của service | Hiển thị trong spans |
| OTEL_TRACES_SAMPLER | Chiến lược sampler | always_on, always_off, parentbased_always_on, ratio:0.x |
| DATABASE_URL | DSN Postgres | Tương thích pgx |
| MEMCACHE_ADDR | Địa chỉ Memcached | host:port |

## Ví dụ sampler
| Giá trị | Ý nghĩa |
|--------|---------|
| always_on | Ghi tất cả span |
| always_off | Tắt tracing |
| parentbased_always_on | Tôn trọng span cha hoặc lấy mẫu root mới |
| ratio:0.2 | Lấy mẫu xác suất 20% |

## Nội suy / Nội bộ
Khi khởi tạo tracer sẽ làm sạch endpoint: nếu có hậu tố `/v1/traces` sẽ bị cắt bỏ.

## Thêm cấu hình mới
1. Thêm field trong struct `Config`.
2. Nạp từ `os.Getenv` với giá trị mặc định hợp lý.
3. Ghi lại tại đây.
