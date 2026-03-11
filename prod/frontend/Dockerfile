FROM node:20-alpine AS builder

WORKDIR /app

COPY prod/frontend/package*.json ./
RUN npm install

COPY prod/frontend/ ./
RUN npm run build

FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY prod/frontend/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
