BINARY_NAME=ed
ROOT_PACKAGE=ed
ROOT_DIR := $(shell realpath .)
GO := go
BINDATA := $(GOPATH)/bin/go-bindata

all: clean format package

gen-doc:
	@swag init -g cmd/$(BINARY_NAME)/main.go -o ./internal/$(BINARY_NAME)/docs

run:
	$(GO) run cmd/$(BINARY_NAME)/main.go

package: bindata
	@GOOS=linux GOARCH=amd64 $(GO) build -tags=jsoniter \
    -o build/$(BINARY_NAME) cmd/$(BINARY_NAME)/main.go

clean:
	@rm -f build/$(BINARY_NAME)
	@$(GO) clean

format:
#	@gofmt -s -w .
#	@$(GOPATH)/bin/gci write ./build/ --section Standard --section Default --section "Prefix(tophant.com/ed)" --skip-generated

bindata:
	$(BINDATA) -o gen/configs.go -pkg gen ./asset/...  ./resources/...
