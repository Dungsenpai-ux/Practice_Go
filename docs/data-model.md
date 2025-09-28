# Mô hình dữ liệu

## Movie
| Trường | Kiểu | Ghi chú |
|--------|------|---------|
| id | int | Khóa chính |
| title | text | Tiêu đề phim |
| year | int | Năm phát hành (tùy chọn) |
| genres | text | Danh sách thể loại phân tách bằng dấu `|` |

SQL (từ migration):
```sql
CREATE TABLE movies (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  year INT,
  genres TEXT
);
CREATE INDEX idx_movies_title ON movies USING GIN (to_tsvector('english', title));
```

## Health Response DTO
| Trường | Kiểu | Mô tả |
|--------|------|------|
| status | string | Luôn là "ok" nếu khỏe mạnh |
| version | string | Lấy từ config / biến môi trường version |
| time | Chuỗi RFC3339 | Thời điểm UTC hiện tại |

## Chiến lược Cache
- Cache dương: lưu toàn bộ JSON phim trong 5 phút.
- Cache âm: lưu `{ "error": "không tìm thấy phim" }` trong 30s khi không tìm thấy.
- Khóa cache: `movie:{id}`.
