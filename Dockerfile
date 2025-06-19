# 构建阶段
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o moe-counter

# 运行阶段
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app

# 复制可执行文件
COPY --from=builder /app/moe-counter .

# 创建数据目录
RUN mkdir /data

# 设置持久化存储
VOLUME /data

# 暴露端口
EXPOSE 8088

# 设置默认环境变量
ENV DB_PATH=/data/data.db \
    PORT=8088

# 使用 start 参数启动
CMD ["./moe-counter", "start", "--db", "/data/data.db", "--port", "8088"]