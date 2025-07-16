FROM nginx:1.25-alpine

# Cài đặt njs module
RUN apk add --no-cache nginx-mod-http-js

# Dọn file default
RUN rm -f /etc/nginx/conf.d/default.conf

# Copy trực tiếp các file cấu hình và script
COPY nginx.conf /etc/nginx/nginx.conf
COPY conf.d/api-gateway.conf /etc/nginx/conf.d/api-gateway.conf
COPY auth.js /etc/nginx/auth.js

EXPOSE 80