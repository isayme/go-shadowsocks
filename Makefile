.PHONY: dev-server
dev-server:
	@CONF_FILE_PATH=${PWD}/config/config.dev.json go run main.go server

.PHONY: dev-local
dev-local:
	@CONF_FILE_PATH=${PWD}/config/config.dev.json go run main.go local
