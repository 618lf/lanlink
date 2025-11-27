# LanLink Makefile

.PHONY: all build clean run test help

# 默认目标
all: build

# 编译当前平台
build:
	@echo "编译 LanLink..."
	@go build -o lanlink

# 编译所有平台
build-all:
	@echo "编译所有平台版本..."
	@mkdir -p dist
	@echo "编译 Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 go build -o dist/lanlink-linux-amd64
	@echo "编译 macOS (arm64)..."
	@GOOS=darwin GOARCH=arm64 go build -o dist/lanlink-mac-arm64
	@echo "编译 macOS (amd64)..."
	@GOOS=darwin GOARCH=amd64 go build -o dist/lanlink-mac-amd64
	@echo "编译 Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 go build -o dist/lanlink-windows-amd64.exe
	@echo "编译完成！"
	@ls -lh dist/

# 运行（需要sudo）
run: build
	@echo "运行 LanLink (需要管理员权限)..."
	@sudo ./lanlink

# 清理
clean:
	@echo "清理编译产物..."
	@rm -f lanlink lanlink.exe
	@rm -rf dist/
	@rm -f lanlink.log
	@rm -f config.json
	@echo "清理完成！"

# 测试
test:
	@echo "运行测试..."
	@go test ./...

# 安装依赖
deps:
	@echo "安装依赖..."
	@go mod download
	@go mod tidy

# 生成配置文件
config:
	@if [ ! -f config.json ]; then \
		echo "生成配置文件..."; \
		cp config.example.json config.json; \
		echo "已创建 config.json"; \
	else \
		echo "config.json 已存在"; \
	fi

# 帮助
help:
	@echo "LanLink Makefile 命令:"
	@echo "  make build       - 编译当前平台"
	@echo "  make build-all   - 编译所有平台"
	@echo "  make run         - 编译并运行（需要sudo）"
	@echo "  make clean       - 清理编译产物"
	@echo "  make test        - 运行测试"
	@echo "  make deps        - 安装依赖"
	@echo "  make config      - 生成配置文件"
	@echo "  make help        - 显示帮助"

