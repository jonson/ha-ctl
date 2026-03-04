VERSION := 0.1.0
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: build build-all test test-integration smoke-test clean

build:
	go build -trimpath -ldflags "$(LDFLAGS)" -o ha-ctl .

build-all:
	GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o ha-ctl-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -trimpath -ldflags "$(LDFLAGS)" -o ha-ctl-linux-arm64 .

test:
	go test ./...

test-integration:
	go test -tags=integration ./...

smoke-test: build
	./scripts/smoke-test.sh

clean:
	rm -f ha-ctl ha-ctl-linux-amd64 ha-ctl-linux-arm64
	rm -rf dist/
