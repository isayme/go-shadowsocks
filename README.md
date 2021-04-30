## Shadowsocks

[![Docker Image Version (latest semver)](https://img.shields.io/docker/v/isayme/shadowsocks?sort=semver&style=flat-square)](https://hub.docker.com/r/isayme/shadowsocks)
![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/isayme/shadowsocks?sort=semver&style=flat-square)
![Docker Pulls](https://img.shields.io/docker/pulls/isayme/shadowsocks?style=flat-square)

A fast and memory efficient shadowsocks server in Go.

## Support Methods

- chacha20-ietf-poly1305
- aes-128-gcm, aes-192-gcm, aes-256-gcm
- aes-128-cfb, aes-192-cfb, aes-256-cfb
- aes-128-ctr, aes-192-ctr, aes-256-ctr
- dec-cfb
- rc4-md5, rc4-md5-6
- chacha20, chacha20-ietf
- cast5-cfb
- bf-cfb

## Dev

### ssserver

> CONF_FILE_PATH=/path/to/config.json go run cmd/ssserver/main.go

### sslocal

> CONF_FILE_PATH=/path/to/config.json go run cmd/sslocal/main.go

## Docker

> docker pull isayme/shadowsocks:latest

### Image

> make image

### Pubulish Tag

> make publish
