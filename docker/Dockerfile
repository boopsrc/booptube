# syntax=docker/dockerfile:1

FROM golang:1.22-bookworm AS build

WORKDIR /src
RUN apt-get update && apt-get install -y --no-install-recommends curl ca-certificates make bash xz-utils \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV FETCH_FFMPEG_PLATFORMS=linux
RUN sed -i 's/\r$//' scripts/*.sh && make fetch-deps \
    && chmod +x assets/ytdlp/linux-amd64/yt-dlp \
       assets/ffmpeg/linux-amd64/ffmpeg \
       assets/ffmpeg/linux-amd64/ffprobe \
    && assets/ytdlp/linux-amd64/yt-dlp --version \
    && assets/ffmpeg/linux-amd64/ffmpeg -version \
    && assets/ffmpeg/linux-amd64/ffprobe -version

ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown

RUN CGO_ENABLED=0 go build -trimpath \
    -ldflags "-s -w -X booptube/buildinfo.Version=${VERSION} -X booptube/buildinfo.Commit=${COMMIT} -X booptube/buildinfo.BuildDate=${DATE}" \
    -o /out/booptube-web ./cmd/web

FROM debian:bookworm-slim AS runtime

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    wget \
    libstdc++6 \
    libgcc-s1 \
    libffi8 \
    libssl3 \
    libz1 \
    libbz2-1.0 \
    liblzma5 \
    libexpat1 \
    && rm -rf /var/lib/apt/lists/* \
    && useradd -r -u 10001 -m -d /app appuser

WORKDIR /app

COPY --from=build /out/booptube-web /app/booptube-web

RUN mkdir -p /data/downloads && chown -R appuser:appuser /data/downloads /app

USER appuser

ENV BOOPTUBE_ADDR=:8080 \
    BOOPTUBE_DOWNLOAD_DIR=/data/downloads \
    BOOPTUBE_MAX_CONCURRENT=4 \
    BOOPTUBE_FILE_TTL=10m \
    GOMEMLIMIT=1800MiB

EXPOSE 8080

VOLUME ["/data/downloads"]

HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --retries=3 \
    CMD wget -q -O /dev/null http://127.0.0.1:8080/health || exit 1

ENTRYPOINT ["/app/booptube-web"]
