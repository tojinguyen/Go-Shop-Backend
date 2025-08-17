import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { SharedArray } from 'k6/data';

const BASE_URL = 'http://localhost:80'; 

const JWT_TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOiIyN2ZjZTQ3YS0yYmM1LTQ3ZWUtOTMwZi1jNjEwODM2YmQ5ZTAiLCJFbWFpbCI6InRvYWluZ3V5ZW5qb2JAZ21haWwuY29tIiwiUm9sZSI6ImN1c3RvbWVyIiwiZXhwIjoxNzU1NTEzNTg1LCJpYXQiOjE3NTU0MjcxODUsImlzcyI6ImdvLXNob3AtcGxhdGZvcm0iLCJuYmYiOjE3NTU0MjcxODV9.vP7LkzVRsSX5Zh7hu5ueqP_5HbD_bf4M_anuxkiyxQo'; 

const productData = [
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

const products = new SharedArray('products', function () {
  return productData;
});


export const options = {
  stages: [
    { duration: '30s', target: 500 }, 
    { duration: '1m', target: 500 },
    { duration: '15s', target: 0 }, 
  ],
  thresholds: {
    'http_req_failed': ['rate<0.01'],      // Tỷ lệ lỗi dưới 1%
    'http_req_duration': ['p(95)<400'],    // 95% request phải dưới 0.1 giây (API ghi thường chậm hơn đọc)
    'http_req_duration': ['p(99)<900'],    // 99% request phải dưới 0.9 giây
    'checks': ['rate>0.99'],
  },
};


export default function () {
  group('Add Item to Cart API', function () {
    const params = {
      headers: {
        'Authorization': `Bearer ${JWT_TOKEN}`,
        'Content-Type': 'application/json',
      },
    };

    const product = products[Math.floor(Math.random() * products.length)];
    const productID = product.id;

    const payload = JSON.stringify({
      product_id: productID,
      quantity: 1, 
    });

    const url = `${BASE_URL}/api/v1/carts/items`;

    const res = http.post(url, payload, params);

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

  sleep(Math.random() * 2 + 1);
}