# Hướng dẫn Cài đặt

## Yêu cầu tiên quyết
- Go 1.25+
- Docker / Docker Compose V2
- PostgreSQL (cục bộ hoặc container)
- Memcached (dùng `docker-compose.yml` cung cấp hoặc dịch vụ bên ngoài)

## Biến môi trường (Tối thiểu)
| Biến | Mô tả | Ví dụ |
|------|-------|-------|
| APP_ENV | Tên môi trường | local |
| OTEL_EXPORTER_OTLP_ENDPOINT | Endpoint OTLP của Tempo (không thêm /v1/traces phía cuối) | http://localhost:4318 |
| OTEL_SERVICE_NAME | Tên service cho traces | practice-go |
| OTEL_TRACES_SAMPLER | Chiến lược sampler (always_on, always_off, parentbased_always_on, ratio:0.2) | ratio:0.5 |
| DATABASE_URL | Chuỗi kết nối Postgres | postgres://user:pass@localhost:5432/app?sslmode=disable |
| MEMCACHE_ADDR | Địa chỉ Memcached | localhost:11211 |

## Khởi động stack Observability
```bash
docker compose -f docker-compose.tracing.yml up -d
```
Grafana: http://localhost:3000 (admin / admin)

## Chạy migration
Sử dụng `golang-migrate` (cài trên máy):
```bash
migrate -path service/migrations -database "$DATABASE_URL" up
```

## Chạy ứng dụng
```bash
go run ./...
```

## Tạo trace thử nghiệm
```
curl http://localhost:8080/healthz
```
Mở Grafana → Explore → Tempo → Run query.

## Dừng stack
```bash
docker compose -f docker-compose.tracing.yml down
```

## Dọn dẹp
```bash
go mod tidy
```

## Khắc phục sự cố
| Vấn đề | Nguyên nhân | Cách khắc phục |
|--------|-------------|----------------|
| Không thấy span trong Grafana | Sai OTLP endpoint | Dùng http://localhost:4318 (không phải /v1/traces) |
| Lỗi kết nối 4318 bị từ chối | Tempo chưa khởi động | Khởi động stack tracing bằng compose |
| Bộ nhớ cao | Sampler luôn bật (always_on) | Đổi sampler sang ratio:0.2 |
