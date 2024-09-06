.PHONY: all module
EXE_NAME := nicksensor

all:
	GOOS=darwin GOARCH=arm64 go build -o bin/darwin-arm64/$(EXE_NAME)
	GOOS=linux GOARCH=arm64 go build -o bin/linux-arm64/$(EXE_NAME)
	GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/$(EXE_NAME)


module:
	rm -rf module.tar.gz
	go build -o bin/$(EXE_NAME)
	tar czf module.tar.gz bin/$(EXE_NAME)
	rm bin/$(EXE_NAME)

clean:
	rm -rf bin module.tar.gz
