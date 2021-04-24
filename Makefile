.PHONY: dev-server
dev-server:
	@CONF_FILE_PATH=${PWD}/config/config.dev.json go run cmd/ssserver/main.go

.PHONY: dev-local
dev-local:
	@CONF_FILE_PATH=${PWD}/config/config.dev.json go run cmd/sslocal/main.go
