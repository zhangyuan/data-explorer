build:
	go build

clean:
	rm -rf data-explorer
	rm -rf bin/data-explorer-*

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
