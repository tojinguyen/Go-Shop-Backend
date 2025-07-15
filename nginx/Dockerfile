FROM nginx:1.25-alpine

# Cài module njs
RUN apk add --no-cache nginx-mod-http-js

# Dọn file default
RUN rm -f /etc/nginx/conf.d/*

# Load module trong nginx.conf (bắt buộc)
COPY nginx.conf /etc/nginx/nginx.conf
COPY conf.d/ /etc/nginx/conf.d/
COPY auth.js /etc/nginx/auth.js

EXPOSE 80
