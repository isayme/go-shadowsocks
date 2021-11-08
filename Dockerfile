FROM golang:1.15.3-alpine as builder
WORKDIR /app

ARG APP_NAME
ENV APP_NAME ${APP_NAME}
ARG APP_VERSION
ENV APP_VERSION ${APP_VERSION}

COPY . .
RUN mkdir -p ./dist  \
  && GO111MODULE=on GOPROXY=https://goproxy.io,direct go mod download \
  && go build -ldflags "-X github.com/isayme/go-shadowsocks/shadowsocks/util.Name=${APP_NAME} \
  -X github.com/isayme/go-shadowsocks/shadowsocks/util.Version=${APP_VERSION}" \
  -o ./dist/shadowsocks main.go

FROM alpine
WORKDIR /app

ARG APP_NAME
ENV APP_NAME ${APP_NAME}
ARG APP_VERSION
ENV APP_VERSION ${APP_VERSION}

# default config file
ENV CONF_FILE_PATH=/etc/shadowsocks.json

COPY config/config.default.json /etc/shadowsocks.json
COPY --from=builder /app/dist/shadowsocks /app/shadowsocks

CMD ["/app/ssserver"]
