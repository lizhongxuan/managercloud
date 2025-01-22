# 前端构建阶段
FROM node:16 as frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend .
RUN npm run build

# 后端构建阶段
FROM golang:1.20 as backend-builder
WORKDIR /app
COPY go.* ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -o middleware-platform cmd/main.go

# 最终运行阶段
FROM nginx:alpine
COPY --from=frontend-builder /app/frontend/dist /usr/share/nginx/html
COPY --from=backend-builder /app/middleware-platform /usr/local/bin/
COPY deploy/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80 