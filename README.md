## Shadowsocks
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
> CONF_FILE_PATH=go /path/to/config.json run cmd/ssserver/main.go

## Docker
> docker pull isayme/shadowsocks:latest

### Image
> make image

### Pubulish Tag
> make publish
