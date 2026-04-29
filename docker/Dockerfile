# 构建阶段
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -o review-bot ./cmd/server/main.go

# 运行阶段
FROM alpine:latest

# 安装 SVN 客户端
RUN apk --no-cache add subversion

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/review-bot .

# 创建 reviews 目录
RUN mkdir -p /app/reviews
VOLUME /app/reviews

EXPOSE 8080

ENTRYPOINT ["/app/review-bot"]
