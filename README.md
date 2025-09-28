# Layered Pattern trong HTTP API

Phần `http_api` áp dụng **Layered Pattern**, chia thành 3 tầng chính:

---

## 1. Controller (Presentation)
- Nhận request, xử lý nghiệp vụ và trả về response.
- Chịu trách nhiệm xác thực và phân quyền.
- Chuyển đổi dữ liệu giữa **DTO** (Data Transfer Objects) và **model**.

## 2. Model
- Chứa data model của service.
- Định nghĩa các entity và mối quan hệ giữa chúng.
- Chứa các method để validate dữ liệu.

## 3. Service
- Tương tác với cơ sở dữ liệu (ví dụ: PostgreSQL, MongoDB, Redis).
- Tích hợp với dịch vụ bên thứ ba (Twilio, email, ...).
- Gọi đến các service khác trong hệ thống microservice.

---

## Cấu trúc thư mục

```plaintext
/config
    config.go           # Cấu hình và biến môi trường
/controller
    main.go             # Định nghĩa router
    [resource].go       # Controller cho mỗi resource
    dto.go              # Data Transfer Objects
/model
    [resource].go       # Model cho mỗi resource
/service
    db/                 # Tương tác với database
        [resource].go   # Repository cho mỗi resource
        db.go           # Khởi tạo kết nối DB
    [external].go       # Tích hợp với dịch vụ bên ngoài
main.go                 # Entry point của ứng dụng
```

---

> **Lưu ý:**  
> - Các tầng được tách biệt rõ ràng giúp dễ dàng bảo trì, mở rộng và kiểm thử.  
> - Tuân thủ mô hình này giúp tăng tính module hóa và quản lý code hiệu quả hơn trong các dự án Go microservice.

---

## Tracing (OpenTelemetry + Tempo + Grafana)

### Chạy stack quan sát

```bash
docker compose -f docker-compose.tracing.yml up -d
```

Mở Grafana: http://localhost:3000  (admin / admin)

### Biến môi trường

| Variable | Default | Mô tả |
|----------|---------|-------|
| OTEL_EXPORTER_OTLP_ENDPOINT | http://localhost:4318 | Tempo OTLP HTTP endpoint |
| OTEL_SERVICE_NAME | practice-go-api | Tên service trong trace |
| OTEL_TRACES_SAMPLER | parentbased_always_on | Chiến lược sampling |

Sampling theo tỷ lệ (10%):
```
OTEL_TRACES_SAMPLER=ratio:0.1
```

### Tích hợp mã nguồn
1. Khởi tạo tracer trong `main.go` (InitTracer + defer shutdown).
2. `controller/main.go` dùng `otelgin` để tạo span mỗi HTTP request.
3. Repository `service/db/movie.repository.go` tạo span thủ công cho truy vấn DB.

### Xem trace
1. Vào Grafana -> Explore -> Chọn Tempo.
2. Filter Service Name = practice-go-api.
3. Chọn trace để xem cây span và timing.

### Mở rộng tương lai
- Thêm logs correlation (trace_id vào log).
- Thêm metrics OTEL (Prometheus / OpenTelemetry Collector).
- Export trace sang Jaeger nếu cần.
