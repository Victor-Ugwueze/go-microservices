

.DEFAULT_GOAL := build

check_install:
	which swagger || G0111MODULE=off go get github.com/go-swagger/go-swagger/cmd/swagger
.PHONY:check_install

swagger:check_install
	G0111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models
.PHONY:swagger

build:
	echo "HI"
.PHONY:build
