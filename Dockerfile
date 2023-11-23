# 指定基础的go编译镜像
FROM golang:alpine as build

# 指定go的环境变量
ENV GOPROXY=https://goproxy.cn \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    HOME=/root

# 指定工作空间目录，会自动cd到这个目录
WORKDIR /build

# 把项目的依赖配置文件拷贝到容器中，并下载依赖
COPY go.mod .
COPY go.sum .
RUN go mod download

# 把项目的其他所有文件拷贝到容器中
COPY . .

# 编译成可执行二进制文件
RUN go build -o app ./cmd/main.go

# 指定新的运行环境，最终的运行会基于这个坏境，使得最终的镜像非常小
FROM scratch as deploy

# 把编译环境中打包好的可执行文件和配置文件拷贝到当前镜像
COPY --from=build /build/app /

CMD ["/app"]
