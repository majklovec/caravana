BINARY_NAME=bin/caravana
PLATFORMS := darwin/amd64 linux/amd64 linux/arm64 windows/amd64

# Add linker flags
linker_flags := "-X caravana/cmd.version=${RELEASE_VERSION}"

build: clean
	mkdir -p bin
	GOARCH=amd64 GOOS=linux go build -v -ldflags=${linker_flags} -o ${BINARY_NAME} .
	chmod +x bin/*

release: clean $(PLATFORMS)

$(PLATFORMS):
	$(eval GOOS=$(word 1, $(subst /, ,$@)))
	$(eval GOARCH=$(word 2, $(subst /, ,$@)))
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -v -ldflags=${linker_flags} -o $(BINARY_NAME)-$(GOOS)-$(GOARCH) .

.PHONY: release $(PLATFORMS)

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm -rf bin

install: build
	sudo cp bin/* /usr/local/bin/

