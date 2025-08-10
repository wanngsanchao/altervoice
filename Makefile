#wangsanchoa writed at 2025-08-04
APP_NAME := altervoice

# 镜像名称和标签（符合 harbor 仓库格式）
HARBOR_REPO := harbor.primelifescience.com.cn/altervoice
IMAGE_NAME := $(HARBOR_REPO)/$(APP_NAME)
IMAGE_TAG ?= latest  # 允许通过环境变量指定标签，默认latest
FULL_IMAGE := $(IMAGE_NAME):$(IMAGE_TAG)

# 目标平台（Linux amd64）
GOOS := linux
GOARCH := amd64

# 编译选项
LDFLAGS := -s -w  # 去除符号表和调试信息，减小二进制体积

# 默认目标
all: build

# 编译当前目录下所有Go文件为指定名称的二进制程序
build:
	@echo "开始编译Golang程序..."
	@if [ -z "$$(ls *.go 2>/dev/null)" ]; then \
		echo "错误：当前目录下未找到.go文件"; \
		exit 1; \
	fi
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -o $(APP_NAME) *.go
	@echo "编译完成: $(APP_NAME)"

# 构建Docker镜像
docker-build: build
	@echo "开始构建Docker镜像: $(FULL_IMAGE)..."
	docker build -t $(FULL_IMAGE) .
	@echo "镜像构建完成: $(FULL_IMAGE)"

# 运行容器（挂载配置文件并映射端口）
docker-run: docker-build
	@echo "检查宿主机配置文件是否存在..."
	@if [ ! -f "/etc/altervoice/altervoice.json" ]; then \
		echo "错误：宿主机配置文件 /etc/altervoice/altervoice.json 不存在"; \
		exit 1; \
	fi
	@echo "停止已存在的容器（如果有）..."
	docker rm -f $(APP_NAME) 2>/dev/null || true
	@echo "启动容器: $(APP_NAME)..."
	docker run -d \
		--name $(APP_NAME) \
		-p 9000:9000 \
		-v /etc/altervoice/altervoice.json:/etc/altervoice/altervoice.json \
		$(FULL_IMAGE)
	@echo "容器已启动，名称: $(APP_NAME)，访问: http://localhost:9000"

# 停止并删除容器
docker-stop:
	@echo "停止并删除容器: $(APP_NAME)..."
	docker rm -f $(APP_NAME) 2>/dev/null || true
	@echo "容器已停止"

# 推送镜像到Harbor仓库
docker-push: docker-build
	@echo "推送镜像到仓库: $(FULL_IMAGE)..."
	docker push $(FULL_IMAGE)
	@echo "镜像推送完成"

# 清理编译产物和本地镜像
clean: docker-stop
	@echo "清理编译产物..."
	rm -f $(APP_NAME)
	@echo "清理本地镜像: $(FULL_IMAGE)..."
	docker rmi $(FULL_IMAGE) 2>/dev/null || true
	@echo "清理完成"

# 显示镜像信息
image-info:
	@echo "镜像信息: $(FULL_IMAGE)"
	docker images $(FULL_IMAGE)

# 显示帮助信息
help:
	@echo "可用命令:"
	@echo "  make                 - 编译当前目录下所有Go文件为$(APP_NAME)"
	@echo "  make build           - 同上"
	@echo "  make docker-build    - 编译程序并构建镜像: $(FULL_IMAGE)"
	@echo "  make docker-run      - 构建镜像并启动容器（挂载配置文件和端口映射）"
	@echo "  make docker-stop     - 停止并删除容器"
	@echo "  make docker-push     - 构建并推送镜像到Harbor仓库"
	@echo "  make clean           - 清理编译产物和本地镜像"
	@echo "  make image-info      - 显示镜像信息"
	@echo "  make help            - 显示帮助信息"
	@echo "  示例: 构建指定标签的镜像并推送"
	@echo "    IMAGE_TAG=v1.0.0 make docker-push"

.PHONY: all build docker-build docker-run docker-stop docker-push clean image-info help
