# 认证授权中心 Makefile

.PHONY: help build run test clean docker-build docker-run install-tools

# 默认目标
help:
	@echo "可用的命令:"
	@echo "  build        构建应用"
	@echo "  run          运行应用"
	@echo "  test         运行测试"
	@echo "  clean        清理构建文件"
	@echo "  docker-build 构建Docker镜像"
	@echo "  docker-run   运行Docker容器"
	@echo "  install-tools 安装开发工具"
	@echo "  swagger      生成Swagger文档"
	@echo "  migrate      运行数据库迁移"

# 构建应用
build:
	@echo "构建应用..."
	go build -o bin/auth-center main.go

# 运行应用
run:
	@echo "运行应用..."
	go run main.go

# 运行测试
test:
	@echo "运行测试..."
	go test -v ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf bin/
	rm -rf docs/

# 安装开发工具
install-tools:
	@echo "安装开发工具..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 生成Swagger文档
swagger:
	@echo "生成Swagger文档..."
	swag init -g main.go -o docs

# 代码检查
lint:
	@echo "运行代码检查..."
	golangci-lint run

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 构建Docker镜像
docker-build:
	@echo "构建Docker镜像..."
	docker build -f docker/Dockerfile -t auth-center .

# 运行Docker容器
docker-run:
	@echo "运行Docker容器..."
	docker run -d \
		-p 8080:8080 \
		-e DB_HOST=localhost \
		-e DB_PASSWORD=password \
		-e REDIS_HOST=localhost \
		--name auth-center \
		auth-center

# 停止Docker容器
docker-stop:
	@echo "停止Docker容器..."
	docker stop auth-center
	docker rm auth-center

# 使用Docker Compose启动
docker-compose-up:
	@echo "使用Docker Compose启动..."
	docker-compose -f docker/docker-compose.yml up -d

# 停止Docker Compose
docker-compose-down:
	@echo "停止Docker Compose..."
	docker-compose -f docker/docker-compose.yml down

# 查看日志
logs:
	@echo "查看日志..."
	docker logs -f auth-center

# 进入容器
shell:
	@echo "进入容器..."
	docker exec -it auth-center /bin/sh

# 数据库迁移
migrate:
	@echo "运行数据库迁移..."
	# 这里可以添加数据库迁移命令

# 开发环境设置
dev-setup:
	@echo "设置开发环境..."
	@if [ ! -f config/app.ini ]; then \
		cp config/app.ini.example config/app.ini; \
		echo "已创建配置文件 config/app.ini，请修改配置"; \
	fi
	@echo "请确保MySQL和Redis服务已启动"

# 完整测试
test-all: fmt lint test
	@echo "所有测试完成"

# 发布构建
release:
	@echo "构建发布版本..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/auth-center main.go
	@echo "发布版本构建完成: bin/auth-center"
