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
config/
    config.go              # Load biến môi trường, build struct Config (không mở kết nối)
controller/
    main.go                # Router helpers (Gin optional), đăng ký routes
    movies.controller.go   # Controller resource Movie (sử dụng DTO nội bộ controller)
    health.controller.go   # Health endpoint
    dto.go                 # Tập trung DTO (chỉ dùng tại layer Controller)
model/
    movie.model.go         # Domain entity Movie
    health.model.go        # (Nếu cần domain-level struct cho health)
service/
    movie.service.go       # Business logic (gọi repository), không phụ thuộc controller hay DTO
    db/
        db.go                # Khởi tạo pool + migrations
        movie.repository.go  # Repository (truy vấn DB cụ thể Movie)
        seed.go              # Seed dữ liệu ban đầu (nhận *pgxpool.Pool)
    external/
        memcache.go          # Khởi tạo Memcached (graceful degrade)
main.go                  # Entry point: load config -> init db -> seed -> init external -> build services & handlers -> start server
```

### Luồng phụ thuộc
Controller -> Service -> Repository(DB) -> Database
Controller -> DTO (chỉ tại controller)
Service chỉ dùng `model` + repository + external (cache) abstraction.
Không có chiều ngược hoặc ngang cấp.

### Nguyên tắc áp dụng
1. Mỗi layer chỉ biết layer bên dưới.
2. DTO không đi xuống Service/Repository.
3. Repository không trả về DTO, chỉ trả về `model`.
4. External integrations (Memcached, v.v.) được gói trong `service/external` và inject vào handler/service.
5. Không còn logic kết nối DB hoặc cache trong package `config`.

### Khởi tạo (main.go)
1. `cfg := config.Load()`
2. `pool, _ := db.Init(cfg)` + defer pool.Close()
3. `db.SeedData(pool)` (idempotent)
4. `cache := external.InitMemcache(cfg.MemcachedAddr)` (có thể nil)
5. `movieService := service.NewMovieService(pool)`
6. Handlers: `movieHandler := controller.NewMovieHandler(movieService, cache)`
7. Đăng ký routes: `controller.RegisterHTTPRoutes(movieHandler)`
8. `http.ListenAndServe("":"+cfg.Port, nil)`

### TODO / Gợi ý mở rộng
- Thêm interface cho cache để test (MockCache).
- Logging middleware & tracing.
- Thay net/http bằng Gin nhất quán (hoặc ngược lại) tránh trộn.
- Thêm layer use-case nếu nghiệp vụ phức tạp.

---

> **Lưu ý:**  
> - Các tầng được tách biệt rõ ràng giúp dễ dàng bảo trì, mở rộng và kiểm thử.  
> - Tuân thủ mô hình này giúp tăng tính module hóa và quản lý code hiệu quả hơn trong các dự án Go microservice.
