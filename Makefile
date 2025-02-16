BIN_NAME = dwmbar
BUILD_DIR = $(HOME)/.dwm/bin
INSTALL_DIR = /usr/local/bin

.PHONY: all build install clean

all: build

build:
	go build -o $(BIN_NAME)
	strip $(BIN_NAME)

install: build
	mkdir -p $(BUILD_DIR)
	cp $(BIN_NAME) $(BUILD_DIR)/$(BIN_NAME)
	cp $(BIN_NAME) $(INSTALL_DIR)/$(BIN_NAME)

clean:
	rm -f $(BIN_NAME)
