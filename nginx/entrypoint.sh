#!/bin/sh
# entrypoint.sh

set -e

# Tìm tất cả các file .template trong /etc/nginx/conf.d
# và thay thế các biến môi trường trong đó.
# Ví dụ: __JWT_SECRET_KEY__ sẽ được thay bằng giá trị của $JWT_SECRET_KEY
for template_file in /etc/nginx/conf.d/*.template; do
  output_file=$(echo "$template_file" | sed 's/\.template$//')
  echo "Processing template file: $template_file -> $output_file"
  envsubst < "$template_file" > "$output_file"
done

echo "Configuration files processed. Starting Nginx."

# Lệnh này sẽ thực thi lệnh CMD mặc định của image Nginx
exec "$@"