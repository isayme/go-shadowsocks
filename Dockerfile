FROM golang:1.11.0-alpine AS builder

RUN wget https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 -q -O $GOPATH/bin/dep && chmod +x $GOPATH/bin/dep
RUN apk update && apk add --no-cache git && rm -rf /var/cache/apk/*

ARG APP_PKG
WORKDIR /go/src/${APP_PKG}

COPY Gopkg.* ./
RUN dep ensure -vendor-only

COPY . .

ARG APP_VERSION
RUN CGO_ENABLED=0 go build -ldflags "-X main.Version=${APP_VERSION}" -o /app/ssserver cmd/ssserver/main.go

FROM alpine
WORKDIR /app

COPY config/config.default.json /etc/shadowsocks.json
COPY --from=builder /app/ssserver /app/ssserver

CMD ["/app/ssserver"]
