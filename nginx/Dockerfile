FROM nginx:alpine

RUN rm -f /etc/nginx/conf.d/*

# Copy cấu hình
COPY nginx.conf /etc/nginx/nginx.conf
COPY conf.d/ /etc/nginx/conf.d/
COPY auth.js /etc/nginx/auth.js

# Expose port
EXPOSE 80
