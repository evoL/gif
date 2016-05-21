GIF_VERSION := $(shell git describe --tags HEAD)

all:
	go build -ldflags="-X github.com/evoL/gif/version.Version=$(GIF_VERSION)"

install:
	go install -ldflags="-X github.com/evoL/gif/version.Version=$(GIF_VERSION)"
