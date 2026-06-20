YTDLP_VERSION ?= 2026.06.09
FFMPEG_VERSION ?= 8.1.1
BUILD_DIR := .build
STAGING_DIR := installer/staging
BINARY := $(BUILD_DIR)/booptube
BINARY_GUI := $(BUILD_DIR)/booptube-gui
BUNDLE_TAGS := -tags bundled

VERSION ?= $(shell cat VERSION 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS = -s -w \
	-X booptube/buildinfo.Version=$(VERSION) \
	-X booptube/buildinfo.Commit=$(COMMIT) \
	-X booptube/buildinfo.BuildDate=$(DATE)

ifeq ($(OS),Windows_NT)
	GUI_LDFLAGS = $(LDFLAGS) -H=windowsgui
	EXE_EXT := .exe
else
	GUI_LDFLAGS = $(LDFLAGS)
	EXE_EXT :=
endif

.PHONY: fetch-ytdlp fetch-ffmpeg fetch-deps \
	build build-gui build-bundled build-gui-bundled \
	stage stage-portable clean \
	package package-portable \
	package-win package-linux package-macos \
	package-portable-win package-portable-linux package-portable-macos

fetch-ytdlp:
	YTDLP_VERSION=$(YTDLP_VERSION) bash scripts/fetch-ytdlp.sh

fetch-ffmpeg:
	FFMPEG_VERSION=$(FFMPEG_VERSION) bash scripts/fetch-ffmpeg.sh

fetch-deps: fetch-ytdlp fetch-ffmpeg

build: fetch-deps
	mkdir -p $(BUILD_DIR)
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY)$(EXE_EXT) ./cmd/cli

build-gui: fetch-deps
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 go build -trimpath -ldflags "$(GUI_LDFLAGS)" -o $(BINARY_GUI)$(EXE_EXT) ./cmd/gui

build-bundled: fetch-deps
	mkdir -p $(BUILD_DIR)
	go build $(BUNDLE_TAGS) -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY)$(EXE_EXT) ./cmd/cli

build-gui-bundled: fetch-deps
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 go build $(BUNDLE_TAGS) -trimpath -ldflags "$(GUI_LDFLAGS)" -o $(BINARY_GUI)$(EXE_EXT) ./cmd/gui

stage: build-bundled build-gui-bundled
	bash scripts/stage.sh bundled

stage-portable: build build-gui
	bash scripts/stage.sh portable

package-portable-win: build build-gui
	bash scripts/package-portable.sh windows

package-portable-linux: build build-gui
	bash scripts/package-portable.sh linux

package-portable-macos: build build-gui
	bash scripts/package-portable.sh macos

package-win: stage
	bash scripts/package.sh windows

package-linux: stage
	bash scripts/package.sh linux

package-macos: stage
	bash scripts/package.sh macos

package-portable:
	bash scripts/package-portable.sh auto

package:
	bash scripts/package.sh auto

clean:
	rm -rf $(BUILD_DIR) $(STAGING_DIR)
