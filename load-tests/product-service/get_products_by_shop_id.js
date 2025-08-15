import http from 'k6/http';
import { check, sleep, group } from 'k6';

const BASE_URL = 'http://localhost:80'; 

const SHOP_ID = '786f4635-b1d7-4437-b86c-0bca872adf5e'; 

const JWT_TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOiIyN2ZjZTQ3YS0yYmM1LTQ3ZWUtOTMwZi1jNjEwODM2YmQ5ZTAiLCJFbWFpbCI6InRvYWluZ3V5ZW5qb2JAZ21haWwuY29tIiwiUm9sZSI6ImN1c3RvbWVyIiwiZXhwIjoxNzU1MzY4NjI4LCJpYXQiOjE3NTUyODIyMjgsImlzcyI6ImdvLXNob3AtcGxhdGZvcm0iLCJuYmYiOjE3NTUyODIyMjh9.oeoFzynj2BZkswl50ouQCzS_A7PBacWIL3lIpSrcTbE'; 


export const options = {
  // stages xác định các giai đoạn tăng và giảm tải
  stages: [
    { duration: '30s', target: 50 }, // Tăng từ 0 lên 50 VUs trong 30 giây
    { duration: '1m', target: 50 },  // Giữ 50 VUs trong 1 phút để duy trì tải
    { duration: '15s', target: 0 },  // Giảm từ 50 về 0 VUs trong 15 giây
  ],

  // thresholds định nghĩa các điều kiện thành công/thất bại của bài test
  thresholds: {
    // Tỷ lệ request lỗi phải dưới 1%
    http_req_failed: ['rate<0.01'], 
    
    // 95% số request phải hoàn thành dưới 500ms
    http_req_duration: ['p(95)<500'], 

    // Tỷ lệ các 'check' thành công phải trên 99%
    checks: ['rate>0.99'],
  },
};

export default function () {
  // Nhóm các request vào một group để dễ đọc report
  group('Get Products by Shop ID API (Query Params)', function () {
    // Chuẩn bị headers cho request, bao gồm cả token xác thực
    const params = {
      headers: {
        'Authorization': `Bearer ${JWT_TOKEN}`,
      },
    };

    // Chuẩn bị URL với các query parameters
    // Random trang để tránh cache và mô phỏng hành vi người dùng tốt hơn
    const page = Math.floor(Math.random() * 10) + 1; // Lấy ngẫu nhiên từ trang 1 đến 10
    const limit = 20;
    const url = `${BASE_URL}/api/v1/products?shop_id=${SHOP_ID}&page=${page}&limit=${limit}`;

    // Gửi request GET với URL đã được xây dựng
    const res = http.get(url, params);

    // Kiểm tra (check) các điều kiện của response
    check(res, {
      'status is 200': (r) => r.status === 200,
      'response body is not empty': (r) => r.body.length > 0,
      'response is valid JSON': (r) => {
        try {
          r.json();
          return true;
        } catch (e) {
          return false;
        }
      },
      'response has "data" array': (r) => Array.isArray(r.json().data),
      'response has "meta" object': (r) => typeof r.json().meta === 'object',
    });
  });

  // Tạm dừng 1 giây giữa các lần lặp để mô phỏng "think time" của người dùng
  sleep(1);
}

// ====================================================================================
// D. HƯỚNG DẪN SỬ DỤNG
// ====================================================================================
/*
### Các bước để chạy bài test:

**Bước 1: Lấy một `shop_id` hợp lệ**
   - Kết nối vào database của `shop-service` (ví dụ: dùng DBeaver, TablePlus, hoặc psql).
   - Chạy câu lệnh SQL sau để lấy một ID ngẫu nhiên:
     ```sql
     SELECT id FROM shops LIMIT 1;
     ```
   - Copy giá trị ID này và dán vào biến `SHOP_ID` ở trên.

**Bước 2: Lấy một `JWT Token` hợp lệ**
   - Kết nối vào database của `user-service`.
   - Lấy email của một user (ví dụ: một `seller`):
     ```sql
     SELECT email FROM user_accounts WHERE user_role = 'seller' LIMIT 1;
     ```
   - Mật khẩu mặc định khi seed là `Password123!`.
   - Sử dụng `curl` hoặc Postman để gửi request đăng nhập và lấy token:
     ```bash
     curl -X POST 'http://localhost:80/api/v1/auth/login' \
     -H 'Content-Type: application/json' \
     --data-raw '{
         "email": "EMAIL_BẠN_VỪA_LẤY_ĐƯỢC",
         "password": "Password123!"
     }' | jq -r '.data.access_token'
     ```
     (Lưu ý: bạn cần cài `jq` để trích xuất token tự động)
   - Copy token nhận được và dán vào biến `JWT_TOKEN` ở trên.

**Bước 3: Chạy k6 test**
   - Mở terminal ở thư mục gốc của dự án.
   - Chạy lệnh:
     ```bash
     k6 run load-tests/product-service/get_products_by_shop_id.js
     ```

**Bước 4: Phân tích kết quả**
   - k6 sẽ hiển thị kết quả chi tiết sau khi chạy xong, bao gồm:
     - Số lượng request thành công/thất bại.
     - Thời gian phản hồi trung bình, p(95), p(99).
     - Kết quả của các `thresholds` bạn đã định nghĩa.
*/