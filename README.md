# Layered Pattern trong HTTP API

Phần `http_api` áp dụng **Layered Pattern**, chia thành 3 tầng chính:

Controller (Presentation):

Nhận request, xử lý nghiệp vụ và trả về response
Chịu trách nhiệm cho việc xác thực và phân quyền
Chuyển đổi dữ liệu giữa DTO (Data Transfer Objects) và model
Model:

Chứa data model của service
Định nghĩa các entity và mối quan hệ giữa chúng
Chứa các method để validate dữ liệu
Service:

Tương tác với cơ sở dữ liệu (ví dụ: PostgreSQL, MongoDB, Redis)
Tích hợp với dịch vụ bên thứ ba (Twilio, email,...)
Gọi đến các service khác trong hệ thống microservice


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