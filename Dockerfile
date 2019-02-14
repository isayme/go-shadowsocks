FROM golang:1.11.5-alpine AS builder

RUN apk update && apk add git

ARG APP_PKG
WORKDIR /go/src/${APP_PKG}

ENV GO111MODULE=on

COPY go.* ./
RUN go mod download

COPY . .

ARG APP_VERSION
RUN CGO_ENABLED=0 go build -ldflags "-X main.Version=${APP_VERSION}" -o /app/ssserver cmd/ssserver/main.go

FROM alpine
WORKDIR /app

COPY config/config.default.json /etc/shadowsocks.json
COPY --from=builder /app/ssserver /app/ssserver

CMD ["/app/ssserver"]
