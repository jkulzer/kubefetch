# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
INSTALL_DIR=/usr/bin

# Binary name
BINARY_NAME=kubefetch

all: clean build

build:
	$(GOBUILD) -ldflags="-s -w -extldflags '-static'" -o $(BINARY_NAME)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

install:
	cp $(BINARY_NAME) $(INSTALL_DIR)

uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)

.PHONY: all build clean
