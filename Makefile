build:
	go build

clean:
	rm -rf data-explorer
	rm -rf bin/data-explorer-*

install:
	go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
	go install github.com/rakyll/gotest@latest
	go install golang.org/x/tools/cmd/deadcode@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1

serve:
	go run main.go serve -c connections.yaml

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/data-explorer_darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o bin/data-explorer_darwin-arm64

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/data-explorer_linux-amd64

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/data-explorer_windows-amd64

build-all: clean build-macos build-linux build-windows

compress-linux:
	upx ./bin/data-explorer_linux*
