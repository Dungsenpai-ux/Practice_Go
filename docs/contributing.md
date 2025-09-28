# Đóng góp

## Nhánh (Branching)
- `main`: ổn định
- Nhánh tính năng: `feature/<tên-ngắn>`

## Kiểu Commit
Khuyến khích dùng Conventional Commits:
```
feat: add movie search filtering by genre
fix: handle cache miss logging panic
chore: remove prometheus metrics
```
(Có thể giữ tiếng Anh cho tiền tố và nội dung ngắn gọn.)

## Pull Request
- Giữ phạm vi rõ ràng & < 400 dòng diff nếu có thể.
- Mô tả nên có: Tại sao (Why) + Thay đổi gì (What changed).

## Phong cách mã
- Chạy gofmt (mặc định) / `go vet` trước khi đẩy.
- Tránh state toàn cục trừ các bộ đếm expvar.
- Giữ hàm ngắn & truyền context.

## Kiểm thử (Tương lai)
Dự kiến bổ sung:
- Unit test (repository với pgx + testcontainer)
- Handler test (httptest + in-memory exporter cho traces)

## Observability khi phát triển
- Dùng Grafana Explore để kiểm tra spans sau thay đổi.
- Nếu không thấy spans: kiểm endpoint, sampler, và biến môi trường.

## Quy trình phát hành (Ví dụ)
1. Tag phiên bản: `git tag -a v0.1.0 -m "initial"`
2. Đẩy tag: `git push --tags`
3. Tạo ghi chú phát hành (release notes).

## Nhãn Issue (Gợi ý)
| Nhãn | Mục đích |
|------|----------|
| bug | Sửa hành vi sai |
| feat | Tính năng mới |
| docs | Thay đổi tài liệu |
| chore | Refactor / dọn dẹp |
| observability | Liên quan telemetry |

## Bảo mật
Không commit secrets. Dùng `.env` (thêm vào .gitignore) hoặc secret manager.

## Liên hệ
Mở issue hoặc tạo thảo luận.
