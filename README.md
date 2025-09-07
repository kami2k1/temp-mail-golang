# Temp Mail (Golang)

Dịch vụ email tạm thời đơn giản viết bằng Go + Gin. Ứng dụng tạo địa chỉ email ngẫu nhiên, kết nối tới máy chủ IMAP để đọc thư, hiển thị trên giao diện web và cung cấp API nhẹ để lấy danh sách thư, xem chi tiết và tải tệp đính kèm.

## Tính năng

- Sinh email ngẫu nhiên: tự động cấp email tạm (VN/EN) dựa trên danh sách domain trong `domain.json`, tránh trùng trong 24h bằng `blacklist.json`.
- Nhận thư IMAP: đọc các hộp thư (All/Junk/Trash), hiển thị nhanh, xem nội dung HTML và ảnh inline, tải tệp đính kèm.
- API đơn giản: không cần đăng ký, xác định danh tính qua cookie JWT tự cấp sẵn.
- Giao diện gọn gàng: SSR (HTML template) + static assets, hỗ trợ dark mode.

## Cấu trúc thư mục

```
cmd/api/main.go          # entrypoint khởi động app
internal/boot            # khởi tạo router, boot server
internal/handlers        # handlers cho trang HTML và API (/messages, /randomize, ...)
internal/middleware      # middleware cấp cookie JWT, gắn email vào context
internal/stmp            # kết nối IMAP, đọc thư, lưu tạm trong RAM, quản lý đính kèm
internal/repository      # lớp truy vấn dữ liệu thư từ bộ nhớ
internal/utils           # sinh email, JWT, blacklist, domain list
internal/database        # khung kết nối DB (chưa sử dụng trong luồng chính)
web/templates            # template HTML (/, /gioi-thieu, /api)
web/static               # CSS/JS và assets tĩnh
domain.json              # danh sách domain để sinh email
blacklist.json           # lưu tạm email đã phát trong 24h
```

## Yêu cầu

- Go 1.21+ (khuyến nghị bản mới gần nhất).
- Máy chủ IMAP (ví dụ Gmail IMAP: cần App Password). Ứng dụng đọc nội dung hộp thư từ IMAP; SMTP hiện chưa dùng trong luồng chính.

## Cấu hình (.env)

Sao chép `example.env` thành `.env` và cập nhật giá trị:

```
HOST=http://127.0.0.1
PORT=1234
APP_ENV=development            # production để bật Gin ReleaseMode

# IMAP / STMP (SMTP chưa dùng, IMAP bắt buộc)
STMP_IMAP=imap.gmail.com
STMP_IMAP_PORT=993
STMP_USER=your-email@gmail.com
STMP_PASS=your-app-password    # App Password nếu dùng Gmail

# JWT
JWT_SECRET=change-me
```

- Chỉnh `domain.json` để thêm/bớt domain phát hành email tạm.
- File `blacklist.json` sẽ được tạo/cập nhật tự động để tránh cấp trùng email trong 24 giờ.

## Chạy dự án

```bash
# 1) Cài dependencies
go mod tidy

# 2) Chạy trực tiếp
go run ./cmd/api

# Hoặc build binary
go build -o tempmail ./cmd/api && ./tempmail
```

Mặc định server chạy tại `http://HOST:PORT`. Truy cập `/` để dùng giao diện web.

## API nhanh

- GET `/messages`: trả về hộp thư hiện tại và danh sách thư (rút gọn, không bao gồm bodyHtml).
- GET `/messages/:uid`: trả về chi tiết một thư (có bodyHtml, danh sách đính kèm).
- POST `/randomize`: tạo địa chỉ email mới và cập nhật cookie danh tính.
- GET `/attachments/:id`: tải tệp đính kèm theo id.

Ghi chú: Server tự quản lý danh tính bằng cookie `token` (JWT, hạn ~24h). Lần đầu gọi API/giao diện sẽ tự tạo email.

## Ghi chú & Hạn chế

- Thư và đính kèm được lưu tạm trong bộ nhớ; khởi động lại sẽ mất.
- Dự án dùng IMAP để nhận thư; không gửi thư ra ngoài.
- Không phù hợp cho dữ liệu nhạy cảm/quan trọng. Chỉ dùng cho mục đích thử nghiệm/đăng ký nhanh.

---

Nếu cần thêm hướng dẫn chi tiết (Docker, deploy, mở rộng API), hãy mở issue hoặc yêu cầu trong dự án.

