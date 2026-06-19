YTDLP_VERSION ?= 2026.06.09
FFMPEG_VERSION ?= 8.1.1
BUILD_DIR := .build
BINARY := $(BUILD_DIR)/booptube
BINARY_GUI := $(BUILD_DIR)/booptube-gui

VERSION ?= $(shell cat VERSION 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS = -s -w \
	-X booptube/buildinfo.Version=$(VERSION) \
	-X booptube/buildinfo.Commit=$(COMMIT) \
	-X booptube/buildinfo.BuildDate=$(DATE)

ifeq ($(OS),Windows_NT)
	GUI_LDFLAGS = $(LDFLAGS) -H=windowsgui
else
	GUI_LDFLAGS = $(LDFLAGS)
endif

.PHONY: fetch-ytdlp fetch-ffmpeg fetch-deps build build-gui clean

fetch-ytdlp:
	YTDLP_VERSION=$(YTDLP_VERSION) ./scripts/fetch-ytdlp.sh

fetch-ffmpeg:
	FFMPEG_VERSION=$(FFMPEG_VERSION) ./scripts/fetch-ffmpeg.sh

fetch-deps: fetch-ytdlp fetch-ffmpeg

build: fetch-deps
	mkdir -p $(BUILD_DIR)
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/cli

build-gui: fetch-deps
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 go build -trimpath -ldflags "$(GUI_LDFLAGS)" -o $(BINARY_GUI) ./cmd/gui

clean:
	rm -rf $(BUILD_DIR)
