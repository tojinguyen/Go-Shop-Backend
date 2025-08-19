import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { SharedArray } from 'k6/data';

// --- CẤU HÌNH ---
const BASE_URL = 'http://localhost:80'; // Địa chỉ API Gateway

// !!! THAY THẾ BẰNG TOKEN HỢP LỆ CỦA BẠN !!!
const JWT_TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOiIyN2ZjZTQ3YS0yYmM1LTQ3ZWUtOTMwZi1jNjEwODM2YmQ5ZTAiLCJFbWFpbCI6InRvYWluZ3V5ZW5qb2JAZ21haWwuY29tIiwiUm9sZSI6ImN1c3RvbWVyIiwiZXhwIjoxNzU1NzA4MjcwLCJpYXQiOjE3NTU2MjE4NzAsImlzcyI6ImdvLXNob3AtcGxhdGZvcm0iLCJuYmYiOjE3NTU2MjE4NzB9.0Oj9mUx5_iT_CBPT1uJuPRZKskzckLSRFfpwTxqU3Kk';

// Dữ liệu mẫu để tạo đơn hàng. Bạn có thể thêm nhiều sản phẩm và shop khác nhau.
const orderData = new SharedArray('order-data', function () {
  return [
    {
      shop_id: '786f4635-b1d7-4437-b86c-0bca872adf5e',
      shipping_address_id: '0d39f076-2511-4380-9c0d-8d2fede60bb2', 
      products: [
        { product_id: 'ca942fb3-12ac-4cf8-9966-dbd034cd3c74', quantity: 1 },
      ],
    },
  ];
});
// --------------------

// --- CẤU HÌNH K6 ---
export const options = {
  // Định nghĩa các giai đoạn của bài test (tăng tải, giữ tải, giảm tải)
  stages: [
    { duration: '30s', target: 50 },    // Tăng dần lên 50 người dùng ảo trong 30 giây
    { duration: '1m', target: 50 },     // Giữ 50 người dùng ảo trong 1 phút
    { duration: '30s', target: 100 },   // Tăng lên 100 người dùng ảo trong 30 giây
    { duration: '1m', target: 100 },    // Giữ 100 người dùng ảo trong 1 phút
    { duration: '30s', target: 0 },     // Giảm về 0 người dùng trong 30 giây
  ],
  // Định nghĩa các ngưỡng performance (pass/fail)
  thresholds: {
    'http_req_failed': ['rate<0.01'],      // Tỷ lệ request lỗi phải dưới 1%
    'http_req_duration': ['p(95)<500'],    // 95% số request phải hoàn thành dưới 500ms
    'http_req_duration': ['p(99)<1500'],   // 99% số request phải hoàn thành dưới 1.5s
    'checks': ['rate>0.99'],               // Tỷ lệ các check thành công phải trên 99%
  },
};
// --------------------

// --- HÀM TEST CHÍNH ---
export default function () {
  group('Create Order API', function () {
    const params = {
      headers: {
        'Authorization': `Bearer ${JWT_TOKEN}`,
        'Content-Type': 'application/json',
      },
    };

    // Chọn một kịch bản đặt hàng ngẫu nhiên từ dữ liệu đã định nghĩa
    const orderPayload = orderData[Math.floor(Math.random() * orderData.length)];
    
    // Gửi request POST
    const res = http.post(
      `${BASE_URL}/api/v1/orders`,
      JSON.stringify({
        shop_id: orderPayload.shop_id,
        shipping_address_id: orderPayload.shipping_address_id,
        items: orderPayload.products,
      }),
      params
    );

    // Kiểm tra kết quả trả về
    check(res, {
      'status is 201 (Created)': (r) => r.status === 201, // API tạo mới thường trả về 201
      'response has "success: true"': (r) => {
        try {
          return r.json().success === true;
        } catch (e) {
          return false; // Lỗi nếu không parse được JSON
        }
      },
      'response has an order ID': (r) => {
        try {
          const data = r.json().data;
          return data && typeof data.id === 'string' && data.id !== '';
        } catch (e) {
          return false;
        }
      },
    });
  });

  // Tạm dừng ngẫu nhiên từ 1-3 giây để mô phỏng hành vi người dùng thật
  sleep(Math.random() * 2 + 1);
}