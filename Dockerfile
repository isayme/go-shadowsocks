FROM alpine
WORKDIR /app

# default config file
ENV CONF_FILE_PATH=/etc/shadowsocks.json

COPY config/config.default.json /etc/shadowsocks.json
COPY ./dist/ssserver /app/ssserver

CMD ["/app/ssserver"]
