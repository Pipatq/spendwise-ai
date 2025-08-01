FROM nginx:1.25-alpine
# ลบ config เริ่มต้นทิ้ง
RUN rm /etc/nginx/conf.d/default.conf
# คัดลอก config ที่เราสร้างขึ้นไปแทน
COPY nginx.conf /etc/nginx/nginx.conf