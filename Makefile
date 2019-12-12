
##
# make                                     テストと全リリースビルド
# make darwin                              darwin用リリースビルド
# make DEBUG=1 darwin                      darwin用デバッグビルド
# make ARGS="main.go -arg1 -arg2=Xyz" run  即時実行
##

NAME := study-golang1
VERSION := v0.0.1
REVISION := $(shell git rev-parse --short HEAD)

##
# ビルドオプション
##

GOARCH := amd64
LDFLAGS := -X 'main.Name=$(NAME)' \
           -X 'main.Version=$(VERSION)' \
           -X 'main.Revision=$(REVISION)' \
           -extldflags '-static'

ifeq ($(DEBUG), 1)
	BUILD_OPTIONS := -race -tags DEBUG -ldflags="$(LDFLAGS)"
	BUILD_MODE := debug
	CGO_ENABLED := 1
else
	BUILD_OPTIONS := -ldflags="-s -w $(LDFLAGS)"
	BUILD_MODE := release
	CGO_ENABLED := 0
endif

##
# 使用するgoコマンドの決定。バージョンとかOS環境とかいろいろあって自動判定はあきらめました。
##

NATIVE_GO := GOPATH= CGO_ENABLED=$(CGO_ENABLED) go
DOCKER_GO := docker run -it -v "$(PWD):/go" -e GOPATH= -e CGO_ENABLED=$(CGO_ENABLED) golang:1.13 go
GO := $(NATIVE_GO)
GO4LINUX := $(DOCKER_GO)
GO4DARWIN := $(NATIVE_GO)
GO4WINDOWS := $(DOCKER_GO)

##
# 作業対象
##

GO_SRCS := $(shell find . -type f -name '*.go')

all: test linux darwin windows

##
# ビルド成果物
##

target/$(BUILD_MODE)/$(NAME)-linux-$(GOARCH): $(GO_SRCS)
	GOOS=linux GOARCH=$(GOARCH) $(GO4LINUX) build $(BUILD_OPTIONS) -o target/$(BUILD_MODE)/$(NAME)-linux-$(GOARCH)

target/$(BUILD_MODE)/$(NAME)-darwin-$(GOARCH): $(GO_SRCS)
	GOOS=darwin GOARCH=$(GOARCH) $(GO4DARWIN) build $(BUILD_OPTIONS) -o target/$(BUILD_MODE)/$(NAME)-darwin-$(GOARCH)

target/$(BUILD_MODE)/$(NAME)-windows-$(GOARCH).exe: $(GO_SRCS)
	GOOS=windows GOARCH=$(GOARCH) $(GO4WINDOWS) build $(BUILD_OPTIONS) -o target/$(BUILD_MODE)/$(NAME)-windows-$(GOARCH).exe

target/tests.xml: $(GO_SRCS)
	$(GO) build -o ./build/bin/go2xunit github.com/tebeka/go2xunit
	$(GO) test -v ./... 2>&1 | ./build/bin/go2xunit -output target/tests.xml
	$(GO) mod tidy # 現状のtidyは実行ファイルへの依存を検知できないためここでgo.modをrevertする

##
# タスク
##

# 実行バイナリファイルを作成します
linux: target/$(BUILD_MODE)/$(NAME)-linux-$(GOARCH)
darwin: target/$(BUILD_MODE)/$(NAME)-darwin-$(GOARCH)
windows: target/$(BUILD_MODE)/$(NAME)-windows-$(GOARCH).exe

# xUnit互換のテスト結果ファイルを作成します
xunit: target/tests.xml

# コード検査を実施します
vet:
	$(GO) vet $(BUILD_OPTIONS) ./...
	$(GO) vet $(BUILD_OPTIONS) ./...

# 単体テストを実施します
test:
	$(GO) mod verify
	$(GO) test -v ./...

# コードフォーマッタを適用します
fmt:
	$(GO) mod tidy
	$(GO) fmt $(shell $(GO) list ./...)

# ビルド生成ファイルを全掃除します
clean:
	-rm -rf target/*
	-rm -rf ./build/*

# 即時実行します
run:
	@ $(GO) run $(BUILD_OPTIONS) $(ARGS)

.PHONY: all linux darwin windows xunit vet test fmt clean run
