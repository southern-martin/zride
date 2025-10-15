
# 🚗 Ứng dụng Kết Nối Xe Chiều Về (Zalo Mini App)

## 1. Tổng Quan Dự Án
Ứng dụng giúp kết nối giữa **tài xế xe khách / taxi / xe dịch vụ đang chạy rỗng chiều về** và **hành khách có nhu cầu đi cùng tuyến đường đó**.  
Hoạt động trên nền tảng **Zalo Mini App**, tận dụng **Zalo OA API** và **AI Matching Engine** để gợi ý chuyến đi phù hợp.

---

## 2. Mục Tiêu MVP
- Tạo phiên bản khả dụng sớm trong 2-3 tháng đầu.
- Cho phép tài xế đăng chuyến chiều về.
- Cho phép hành khách tìm kiếm hoặc nhận gợi ý chuyến đi phù hợp.
- Tích hợp thanh toán cơ bản (chuyển khoản hoặc ZaloPay).

---

## 3. Thành Phần Chính

### 3.1 Frontend (Zalo Mini App)
- Hiển thị danh sách chuyến xe rỗng.
- Cho phép hành khách đăng yêu cầu đi nhờ.
- Giao diện thân thiện, tối ưu cho mobile.
- Tích hợp chat với tài xế thông qua Zalo OA.

### 3.2 Backend Services (Microservices Architecture)
- **API Gateway**: định tuyến và xác thực yêu cầu.
- **Auth Service**: xác thực người dùng (Zalo OAuth).
- **Trip Service**: quản lý chuyến xe, điểm đi, điểm đến, trạng thái.
- **User Service**: lưu thông tin tài xế và hành khách.
- **AI Matching Service**: sử dụng gợi ý tuyến đường, giá vé, hoặc hành khách phù hợp.
- **Payment Service**: xử lý thanh toán (ZaloPay hoặc ví điện tử khác).

### 3.3 Database
- PostgreSQL cho dữ liệu giao dịch chính.
- Redis cho cache và xử lý matching nhanh.
- MongoDB (tuỳ chọn) để lưu logs hoặc dữ liệu AI training.

---

## 4. Luồng Dữ Liệu
1. Tài xế đăng nhập qua Zalo OA.
2. Đăng chuyến: điểm đi, điểm đến, thời gian, số ghế trống, giá dự kiến.
3. AI Matching Service đề xuất cho hành khách phù hợp.
4. Hành khách chọn chuyến, xác nhận, thanh toán.
5. Chat hoặc gọi điện trực tiếp qua Zalo OA API.

---

## 5. AI Matching Logic
- Gợi ý chuyến dựa trên:
  - Khoảng cách địa lý (sử dụng Google Maps API hoặc OpenStreetMap).
  - Thời gian khởi hành gần nhất.
  - Xếp hạng tài xế và phản hồi của khách trước đó.
- Mô hình gợi ý: dùng **cosine similarity** hoặc **kNN model** trong giai đoạn MVP.

---

## 6. Kiến Trúc Hệ Thống

### Sơ Đồ (Mô Tả)
- **Zalo Mini App (Frontend)**  
  ⬇️ Gửi request đến  
- **API Gateway**  
  ⬇️ Phân phối đến  
  - Auth Service  
  - Trip Service  
  - User Service  
  - Payment Service  
  - AI Matching Service  
  ⬇️ Lưu dữ liệu tại  
  - PostgreSQL / Redis / MongoDB  

Triển khai trên cloud (AWS/GCP) hoặc VPS tại Việt Nam (VDI, ViettelCloud).

---

## 7. Triển Khai & Công Nghệ
- **Ngôn ngữ**: Golang (backend), ReactJS (mini app), Python (AI service)
- **Frameworks**: Gin/Fiber, Zalo Mini App SDK, FastAPI (cho AI)
- **Cơ sở dữ liệu**: PostgreSQL, Redis
- **Triển khai**: Docker + Traefik + CI/CD (GitHub Actions)
- **Tích hợp**: Zalo OA API, ZaloPay API

---

## 8. Lộ Trình Phát Triển
| Giai đoạn | Thời gian | Nội dung chính |
|------------|------------|----------------|
| **P1 - Nghiên cứu** | Tuần 1-2 | Xác định yêu cầu, nghiên cứu Zalo API |
| **P2 - Thiết kế hệ thống** | Tuần 3-4 | Thiết kế kiến trúc, database, API |
| **P3 - Phát triển MVP** | Tháng 2-3 | Xây dựng Mini App + Backend cơ bản |
| **P4 - Tích hợp AI Matching** | Tháng 4 | Gợi ý tuyến và hành khách tự động |
| **P5 - Thử nghiệm thực tế** | Tháng 5 | Thử nghiệm tại một khu vực cụ thể (Bình Phước – Sài Gòn) |
| **P6 - Ra mắt & Marketing** | Tháng 6 | Tích hợp thanh toán, quảng bá qua OA & nhóm Zalo |

---

## 9. Kế Hoạch Doanh Thu
- **Phí dịch vụ**: 10–15% trên mỗi giao dịch thành công.
- **Quảng cáo**: hiển thị banner hoặc gợi ý dịch vụ cho tài xế (bảo dưỡng, nghỉ ngơi...).
- **Gói đăng ký tài xế chuyên nghiệp**: ưu tiên hiển thị & gợi ý khách.

---

## 10. Tiềm Năng Phát Triển
- Mở rộng sang tuyến cố định, logistics hàng hóa.
- Tích hợp định vị GPS real-time.
- Xây dựng app riêng (Android/iOS) sau khi Zalo Mini App ổn định.
- Ứng dụng mô hình AI nâng cao: dự đoán giá, tối ưu tuyến, xếp lịch tự động.

---

## 11. Kết Luận
Dự án này tận dụng tốt nền tảng Zalo – nơi có lượng người dùng lớn tại Việt Nam, chi phí phát triển thấp, và khả năng lan tỏa nhanh chóng.  
Với MVP được xây dựng đúng hướng, đây có thể là **bước khởi đầu khả thi** cho một hệ sinh thái dịch vụ vận chuyển thông minh.

---

🧭 **Người phụ trách:** Nguyễn Phạm Văn Tân  
📅 **Ngày cập nhật:** 15/10/2025  
🧩 **Trạng thái:** MVP Planning
