NAME := k8s-unused-secret-detector

SRCS     := $(shell find . -type f -name '*.go')
LDFLAGS  := -ldflags="-s -w -extldflags \"-static\""

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS)
	GO111MODULE=on go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: clean
clean:
	rm -rf bin/*

.PHONY: test
test:
	GO111MODULE=on go test ./...
