# Kiến trúc

Tiếp cận phân lớp với ranh giới package rõ ràng.

```
+------------------+
|   Controller     |  -> Handler Gin, định tuyến HTTP
+---------+--------+
          |
+---------v--------+
|     Service      |  -> Điều phối repo, cache, validation
+---------+--------+
          |
+---------v--------+
|   Repository     |  -> Truy vấn PGX tới PostgreSQL
+---------+--------+
          |
+---------v--------+
|   Database       |
+------------------+
```

## Luồng Observability
```
Request -> Gin (middleware otelgin) -> Handler -> Repository call
  |                                             |
  +-- root span ------------------------------- +
         child spans (truy vấn DB, tra cache)
```

## Mẫu Cache
- Read-through kèm negative caching.
- Không có danh sách vô hiệu hóa tổng quát: POST chỉ xóa khóa đơn lẻ.

## Lựa chọn thiết kế
| Khu vực | Quyết định | Lý do |
|---------|------------|-------|
| Framework | Gin | Trưởng thành, tích hợp OTEL tốt |
| Tracing | Gửi OTLP trực tiếp tới Tempo | Đơn giản (không qua collector) |
| Metrics | Chỉ dùng expvar | Nhẹ sau khi bỏ Prometheus |
| Tầng DB | Hàm repository đơn giản | Giữ service mỏng |
| Cache | Memcached | Minh họa phụ thuộc ngoài |

## Đánh đổi
- Chưa có abstraction domain service (giữ tối giản).
- Không chạy migrations trong app (trách nhiệm bên ngoài).
- Thông báo lỗi một phần bằng tiếng Việt cho phù hợp bối cảnh local.

## Nâng cấp tương lai
- Thêm timeout context cho mỗi lời gọi DB.
- Thêm endpoint nạp hàng loạt.
- Thêm circuit breaker cho DB/cache.
- Thêm logging có cấu trúc & correlation ID.
