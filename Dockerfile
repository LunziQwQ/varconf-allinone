# 构建阶段
FROM golang:1.23 AS build

# 设置工作目录
WORKDIR /varconf/

# 下载后端代码并解压
COPY . .

# 编译静态应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -mod=vendor .

# 生产阶段
FROM alpine:latest

RUN apk add --no-cache --update ca-certificates tzdata \
    && update-ca-certificates 2>/dev/null || true
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

# 设置工作目录
WORKDIR /varconf/

# 从build阶段拷贝二进制文件
COPY --from=build /varconf/varconf .
COPY --from=build /varconf/varconf-ui/ ./varconf-ui/

# 添加到环境变量
ENV PATH /varconf:$PATH

# 启动命令
ENTRYPOINT ["varconf"]
CMD ["--help"]