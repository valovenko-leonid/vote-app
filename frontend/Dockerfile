# build stage
FROM node:22-alpine AS builder
WORKDIR /app
COPY frontend/package.json frontend/vite.config.js ./
RUN npm install
COPY frontend/src ./src
RUN npm run build

# tiny Nginx stage
FROM nginx:1.27-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY frontend/nginx.conf /etc/nginx/nginx.conf
EXPOSE 80