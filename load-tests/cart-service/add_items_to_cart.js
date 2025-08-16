// File: load-tests/cart-service/add_item_to_cart.js

import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { SharedArray } from 'k6/data';

// ====================================================================================
// A. CẤU HÌNH BÀI TEST - Vui lòng cập nhật các giá trị này
// ====================================================================================

const BASE_URL = 'http://localhost:80'; 

// JWT Token hợp lệ. Hãy lấy token mới bằng cách gọi API login.
const JWT_TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOiIyN2ZjZTQ3YS0yYmM1LTQ3ZWUtOTMwZi1jNjEwODM2YmQ5ZTAiLCJFbWFpbCI6InRvYWluZ3V5ZW5qb2JAZ21haWwuY29tIiwiUm9sZSI6ImN1c3RvbWVyIiwiZXhwIjoxNzU1MzY5ODE4LCJpYXQiOjE3NTUyODM0MTgsImlzcyI6ImdvLXNob3AtcGxhdGZvcm0iLCJuYmYiOjE3NTUyODM0MTh9.MZCCE_eBgfwcMnfWXIxYchDFSrD1zupvyxrESlhyb_c'; 

// [QUAN TRỌNG] Danh sách các Product ID để thêm vào giỏ hàng.
// Lấy các ID này từ database của product-service. Càng nhiều ID càng tốt.
const productData = [
  // Dán các ID sản phẩm của bạn vào đây
  { id: "ca942fb3-12ac-4cf8-9966-dbd034cd3c74" },
  { id: "e5af5ca8-5150-49bf-b0c3-6092948bc7f1" },
  { id: "bf9a6c04-0891-4294-be3c-440d906fa61d" },
  { id: "77404724-ea66-4593-9e5d-57a4959bb0c4" },
  { id: "133ba97b-e8f1-4f54-abe5-739eecca55cb" },
  { id: "c7a29a9b-f9a2-43ab-a23d-d1a57519fc5b" },
  { id: "1cfec3b7-6b83-4bf9-9475-cf52630ad4b7" },
  { id: "b197c044-b82b-4125-b769-499c3fde63e0" },
  { id: "54b3501f-09d9-4396-85c1-1cd380d4eaeb" },
  { id: "163111e8-9aad-46fe-a7a1-fca1b653552b" },
];

// Sử dụng SharedArray để chia sẻ dữ liệu product IDs giữa các VUs một cách hiệu quả
const products = new SharedArray('products', function () {
  return productData;
});


// ====================================================================================
// B. KỊCH BẢN TẢI (LOAD SCENARIO)
// ====================================================================================
export const options = {
  stages: [
    { duration: '30s', target: 20 }, // Tăng lên 20 người dùng trong 30s
    { duration: '1m', target: 20 },  // Giữ 20 người dùng trong 1 phút
    { duration: '15s', target: 0 },  // Giảm về 0
  ],
  thresholds: {
    'http_req_failed': ['rate<0.01'],      // Tỷ lệ lỗi dưới 1%
    'http_req_duration': ['p(95)<1000'],    // 95% request phải dưới 1 giây (API ghi thường chậm hơn đọc)
    'checks': ['rate>0.99'],
  },
};


export default function () {
  group('Add Item to Cart API', function () {
    // Chuẩn bị headers cho request
    const params = {
      headers: {
        'Authorization': `Bearer ${JWT_TOKEN}`,
        'Content-Type': 'application/json',
      },
    };

    // Chọn ngẫu nhiên một sản phẩm từ danh sách đã chuẩn bị
    const product = products[Math.floor(Math.random() * products.length)];
    const productID = product.id;

    // Tạo body cho request POST
    const payload = JSON.stringify({
      product_id: productID,
      quantity: 1, // Mô phỏng người dùng thêm 1 sản phẩm mỗi lần
    });

    const url = `${BASE_URL}/api/v1/carts/items`;

    // Gửi request POST
    const res = http.post(url, payload, params);

    // Kiểm tra kết quả trả về
    check(res, {
      'status is 200 (OK)': (r) => r.status === 200,
      'response has "success: true"': (r) => {
        try {
          return r.json().success === true;
        } catch (e) {
          return false;
        }
      },
    });
  });

  // Tạm dừng từ 1 đến 3 giây để mô phỏng hành vi người dùng thực tế hơn
  sleep(Math.random() * 2 + 1);
}