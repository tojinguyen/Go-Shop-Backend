import http from 'k6/http';
import { check, sleep, group } from 'k6';

const BASE_URL = 'http://localhost:80'; 

const SHOP_ID = '786f4635-b1d7-4437-b86c-0bca872adf5e'; 

const JWT_TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOiIyN2ZjZTQ3YS0yYmM1LTQ3ZWUtOTMwZi1jNjEwODM2YmQ5ZTAiLCJFbWFpbCI6InRvYWluZ3V5ZW5qb2JAZ21haWwuY29tIiwiUm9sZSI6ImN1c3RvbWVyIiwiZXhwIjoxNzU1MzY5ODE4LCJpYXQiOjE3NTUyODM0MTgsImlzcyI6ImdvLXNob3AtcGxhdGZvcm0iLCJuYmYiOjE3NTUyODM0MTh9.MZCCE_eBgfwcMnfWXIxYchDFSrD1zupvyxrESlhyb_c'; 


export const options = {
  stages: [
    { duration: '30s', target: 100 },  
    { duration: '1m', target: 100 }, 
    { duration: '30s', target: 200 },
    { duration: '1m', target: 200 },  
    { duration: '30s', target: 250 },  
    { duration: '1m', target: 250 },   
    { duration: '30s', target: 0 },   
  ],
  thresholds: {
    'http_req_failed': ['rate<0.02'],     
    'http_req_duration': ['p(95)<800'],
    'checks': ['rate>0.98'],
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
}
