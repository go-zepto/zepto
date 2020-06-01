GOPATH:=$(shell go env GOPATH)

.PHONY: test
test:
	ZEPTO_ENV=test go test -v ./... -cover
