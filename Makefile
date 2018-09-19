APP_NAME := shadowsocks
APP_VERSION := $(shell git describe --tags)
APP_PKG := $(shell echo ${PWD} | sed -e "s\#${GOPATH}/src/\#\#g")

.PHONY: image
image:
	docker build \
	--build-arg APP_PKG=${APP_PKG} \
	--build-arg APP_VERSION=${APP_VERSION} \
	-t ${APP_NAME}:${APP_VERSION} .