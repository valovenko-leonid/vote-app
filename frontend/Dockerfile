# build stage
FROM node:22-alpine AS builder
WORKDIR /app
COPY package.json vite.config.js ./
RUN npm install
COPY . .
RUN npm run build

# tiny Nginx stage
FROM nginx:1.27-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80