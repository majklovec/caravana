BINARY_NAME=bin/caravana

# Get version from git hash
calver := $(shell date +'%y.%m.%d-%H%M%S')
# Add linker flags
linker_flags := "-X caravana/cmd.version=${calver}"

build: clean
	mkdir -p bin
	GOARCH=arm64 GOOS=linux go build -v -ldflags=${linker_flags} -o ${BINARY_NAME}-arm64 .
# GOARCH=amd64 GOOS=darwin go build -v -ldflags=${linker_flags} -o ${BINARY_NAME}-darwin .
	GOARCH=amd64 GOOS=linux go build -v -ldflags=${linker_flags} -o ${BINARY_NAME} .
# GOARCH=amd64 GOOS=windows go build -v -ldflags=${linker_flags} -o ${BINARY_NAME}.exe .
	chmod +x bin/*

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm -rf bin

install: build
	sudo cp bin/* /usr/local/bin