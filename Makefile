all:
	rm -rf bin
	GOOS=darwin GOARCH=arm64 go build -o bin/darwin-arm64/nicksensor
	GOOS=linux GOARCH=arm64 go build -o bin/linux-arm64/nicksensor
	GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/nicksensor
