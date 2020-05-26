CWD := ${CURDIR}
CODE_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BIN_DIR=$(CWD)/build

LAMBDA_FUNC = \
	$(BIN_DIR)/Handler

lambda: $(LAMBDA_FUNC)

$(BIN_DIR)/Handler: $(CODE_DIR)/lambda/Handler/*.go
	cd $(CODE_DIR) && env GOARCH=amd64 GOOS=linux go build -o $(BIN_DIR)/Handler $(CODE_DIR)/lambda/Handler && cd $(CWD)
