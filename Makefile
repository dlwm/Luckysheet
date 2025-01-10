# 获取当前工作目录
CURRENT_DIR := $(shell cd)

# 目标文件和链接路径
TARGET_DIR := $(CURRENT_DIR)\dist
LINK_NAME := $(CURRENT_DIR)\client\public

# 默认目标
all: $(TARGET_FILE)
	@if not exist "$(LINK_NAME)" ( \
		echo Creating symbolic link: $(LINK_NAME) -> $(TARGET_FILE); \
		mklink "$(LINK_NAME)" "$(TARGET_FILE)"; \
	) else ( \
		echo Symbolic link already exists: $(LINK_NAME); \
	)

# 创建目标文件（示例）
$(TARGET_FILE):
	@echo This is a target file. > "$(TARGET_FILE)"

# 清理
clean:
	@if exist "$(TARGET_FILE)" del "$(TARGET_FILE)"
	@if exist "$(LINK_NAME)" del "$(LINK_NAME)"
