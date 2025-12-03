APP_NAME := gh-pr-formatter
BUILD_DIR := dist
VERSION ?= v0.0.0-dev

.PHONY: default build build-all clean

default: build

build:
	@echo "🚀 Building for your current OS (version: $(VERSION))"
	go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME) ./cmd/gh-pr-formatter

build-all:
	@echo "📦 Building all platforms (version: $(VERSION))"

	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME)_mac_arm64 ./cmd/gh-pr-formatter

	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME)_mac_amd64 ./cmd/gh-pr-formatter

	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME)_linux_amd64 ./cmd/gh-pr-formatter

	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME).exe ./cmd/gh-pr-formatter

	@echo "--------------------------------------"
	@echo " 🎉 All binaries built successfully!"
	@echo "--------------------------------------"

clean:
	@rm -rf $(BUILD_DIR)
	@echo "🧹 Cleaned build directory"
