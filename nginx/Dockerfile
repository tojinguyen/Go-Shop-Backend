FROM nginx:1.25-alpine

# Cài đặt các gói cần thiết, bao gồm njs và gettext để có envsubst
RUN apk add --no-cache nginx-mod-http-js gettext

# Dọn file default
RUN rm -f /etc/nginx/conf.d/default.conf

# Copy các file cấu hình và script
COPY nginx.conf /etc/nginx/nginx.conf
# Chú ý: Copy file cấu hình gốc vào một file tạm
COPY conf.d/api-gateway.conf /etc/nginx/conf.d/api-gateway.conf.template
COPY auth.js /etc/nginx/auth.js

# Copy và cấp quyền thực thi cho entrypoint script
COPY entrypoint.sh /docker-entrypoint.d/20-envsubst-on-templates.sh
RUN chmod +x /docker-entrypoint.d/20-envsubst-on-templates.sh

EXPOSE 80