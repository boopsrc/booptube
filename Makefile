YTDLP_VERSION ?= 2026.06.09
YTDLP_BASE := https://github.com/yt-dlp/yt-dlp/releases/download/$(YTDLP_VERSION)
BUILD_DIR := .build
BINARY := $(BUILD_DIR)/booptube

.PHONY: fetch-ytdlp build clean

fetch-ytdlp:
	mkdir -p assets/ytdlp/windows-amd64 assets/ytdlp/linux-amd64 assets/ytdlp/darwin-arm64
	curl -fsSL -o assets/ytdlp/windows-amd64/yt-dlp.exe "$(YTDLP_BASE)/yt-dlp.exe"
	curl -fsSL -o assets/ytdlp/linux-amd64/yt-dlp "$(YTDLP_BASE)/yt-dlp"
	curl -fsSL -o assets/ytdlp/darwin-arm64/yt-dlp "$(YTDLP_BASE)/yt-dlp"
	chmod +x assets/ytdlp/linux-amd64/yt-dlp assets/ytdlp/darwin-arm64/yt-dlp

build: fetch-ytdlp
	mkdir -p $(BUILD_DIR)
	go build -o $(BINARY) .

clean:
	rm -rf $(BUILD_DIR)
