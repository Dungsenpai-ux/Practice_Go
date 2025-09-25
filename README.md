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
