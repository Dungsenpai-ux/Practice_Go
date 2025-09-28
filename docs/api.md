# Tham chiếu API

Base URL: `http://localhost:8080`

## Health
### GET /healthz
Trả về trạng thái dịch vụ và phiên bản.

Phản hồi 200:
```json
{
  "status": "ok",
  "version": "1.0.0",
  "time": "2025-09-28T12:34:56Z"
}
```

## Movies
### POST /movies
Tạo mới một bộ phim.

Request JSON:
```json
{
  "title": "Interstellar",
  "year": 2014,
  "genres": "Adventure|Drama|Sci-Fi"
}
```
Phản hồi 201:
```json
{"id": 123}
```
Lỗi: 400 (body không hợp lệ), 500 (thêm thất bại).

### GET /movies/{id}
Lấy phim theo id. Có cache trong Memcached (cache dương & âm).

Phản hồi 200:
```json
{
  "id": 123,
  "title": "Interstellar",
  "year": 2014,
  "genres": "Adventure|Drama|Sci-Fi"
}
```
Phản hồi 404:
```
không tìm thấy phim
```

### GET /movies/search?q=term&year=YYYY
Tìm phim theo một phần tiêu đề; tham số year tùy chọn (chính xác).

Phản hồi 200:
```json
[
  {"id":1, "title":"Alien", "year":1979, "genres":"Horror|Sci-Fi"},
  {"id":2, "title":"Aliens", "year":1986, "genres":"Action|Sci-Fi"}
]
```

## Xử lý lỗi
Lỗi trả về dạng text thuần với status code phù hợp.

## Tracing
Mỗi endpoint được bọc bởi Gin + otelgin tạo spans; các lời gọi repository phim sinh thêm child spans.

## Giới hạn tốc độ / Auth
Chưa triển khai; có thể thêm gateway hoặc middleware sau.
