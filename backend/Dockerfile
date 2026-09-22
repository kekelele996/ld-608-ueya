FROM golang:1.22-alpine AS build
WORKDIR /app
# 先复制依赖清单，利用 Docker 层缓存
COPY go.mod go.sum ./
RUN go mod download
# 再复制业务源码（src 分层目录 + main.go）
COPY . .
# SQLite 驱动为纯 Go（glebarez/sqlite），关闭 CGO 即可在 alpine 静态构建
ENV CGO_ENABLED=0 GOOS=linux
RUN go mod tidy && go build -ldflags="-s -w" -o server .

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache wget ca-certificates && adduser -D -u 10001 appuser
COPY --from=build /app/server ./server
USER appuser
EXPOSE 3000
CMD ["./server"]
