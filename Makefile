APP_NAME := shadowsocks
APP_VERSION := $(shell git describe --tags)
APP_PKG := $(shell echo ${PWD} | sed -e "s\#${GOPATH}/src/\#\#g")
BUILD_TIME := $(shell date -u +"%FT%TZ")
GIT_REVISION := $(shell git rev-parse HEAD)

.PHONY: dev
dev:
	@CONF_FILE_PATH=${PWD}/config/config.dev.json go run cmd/ssserver/main.go

.PHONY: build
build:
	@go build -ldflags "-X ${APP_PKG}/shadowsocks/util.Name=${APP_NAME} \
	-X ${APP_PKG}/shadowsocks/util.Version=${APP_VERSION} \
	-X ${APP_PKG}/shadowsocks/util.BuildTime=${BUILD_TIME} \
	-X ${APP_PKG}/shadowsocks/util.GitRevision=${GIT_REVISION}" \
	-o ./dist/ssserver cmd/ssserver/main.go

.PHONY: image
image:
	docker build --rm -t ${APP_NAME}:${APP_VERSION} .

.PHONY: publish
publish: image
	docker tag ${APP_NAME}:${APP_VERSION} isayme/${APP_NAME}:${APP_VERSION}
	docker push isayme/${APP_NAME}:${APP_VERSION}
	docker tag ${APP_NAME}:${APP_VERSION} isayme/${APP_NAME}:latest
	docker push isayme/${APP_NAME}:latest
