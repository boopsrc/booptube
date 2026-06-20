#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

version="${FFMPEG_VERSION:-8.1.1}"
minor="${version%.*}"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

mkdir -p assets/ffmpeg/windows-amd64 assets/ffmpeg/linux-amd64 assets/ffmpeg/darwin-arm64

#region agent log
log_debug() {
	local msg="$1" hid="$2" data="$3"
	local ts=$(($(date +%s) * 1000))
	printf '{"sessionId":"bd5798","runId":"verify","hypothesisId":"%s","location":"fetch-ffmpeg.sh","message":"%s","data":%s,"timestamp":%s}\n' \
		"$hid" "$msg" "$data" "$ts" >> "${root}/debug-bd5798.log" 2>/dev/null || true
}
#endregion agent log

need_platforms() {
	if [[ "${FETCH_FFMPEG_PLATFORMS:-}" == "all" ]]; then
		echo "windows linux macos"
		return
	fi
	if [[ -n "${FETCH_FFMPEG_PLATFORMS:-}" ]]; then
		echo "${FETCH_FFMPEG_PLATFORMS//,/ }"
		return
	fi
	case "$(uname -s)" in
		Linux) echo "linux" ;;
		Darwin) echo "macos" ;;
		MINGW*|MSYS*|CYGWIN*|Windows*) echo "windows" ;;
		*) echo "windows linux macos" ;;
	esac
}

curl_download() {
	local out="$1" platform="$2"
	shift 2
	local url code
	for url in "$@"; do
		#region agent log
		log_debug "curl attempt" "C" "{\"platform\":\"${platform}\",\"url\":\"${url}\"}"
		#endregion agent log
		code=$(curl -sSL -w "%{http_code}" -o "$out" \
			-H "User-Agent: booptube-build/1.0 (+https://github.com/booptube/booptube)" \
			"$url" || echo "000")
		if [[ "$code" == "200" && -s "$out" ]]; then
			#region agent log
			log_debug "curl ok" "C" "{\"platform\":\"${platform}\",\"url\":\"${url}\",\"code\":${code}}"
			#endregion agent log
			return 0
		fi
		#region agent log
		log_debug "curl failed" "C" "{\"platform\":\"${platform}\",\"url\":\"${url}\",\"code\":\"${code}\"}"
		#endregion agent log
		rm -f "$out"
	done
	echo "failed to download ffmpeg for ${platform}" >&2
	return 1
}

fetch_windows() {
	local win_zip="${tmpdir}/ffmpeg-win.zip" win_dir="${tmpdir}/ffmpeg-win"
	curl_download "$win_zip" "windows" \
		"https://www.gyan.dev/ffmpeg/builds/packages/ffmpeg-${version}-essentials_build.zip" \
		"https://github.com/GyanD/codexffmpeg/releases/download/${version}/ffmpeg-${version}-essentials_build.zip"
	mkdir -p "$win_dir"
	unzip -q "$win_zip" -d "$win_dir"
	find "$win_dir" -type f \( -name ffmpeg.exe -o -name ffprobe.exe \) -exec cp -f {} assets/ffmpeg/windows-amd64/ \;
}

fetch_linux() {
	local linux_tar="${tmpdir}/ffmpeg-linux.tar.xz" linux_dir="${tmpdir}/ffmpeg-linux"
	curl_download "$linux_tar" "linux" \
		"https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz" \
		"https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-n${minor}-latest-linux64-gpl-${minor}.tar.xz"
	mkdir -p "$linux_dir"
	tar -xf "$linux_tar" -C "$linux_dir"
	find "$linux_dir" -type f \( -name ffmpeg -o -name ffprobe \) ! -path "*/model/*" -exec cp -f {} assets/ffmpeg/linux-amd64/ \;
	chmod +x assets/ffmpeg/linux-amd64/ffmpeg assets/ffmpeg/linux-amd64/ffprobe
}

fetch_macos() {
	local mac_ffmpeg_zip="${tmpdir}/ffmpeg-mac.zip" mac_ffprobe_zip="${tmpdir}/ffprobe-mac.zip" mac_dir="${tmpdir}/ffmpeg-mac"
	for i in 1 2 3; do
		if curl_download "$mac_ffmpeg_zip" "macos-ffmpeg" "https://evermeet.cx/ffmpeg/getrelease/ffmpeg/zip" \
			&& curl_download "$mac_ffprobe_zip" "macos-ffprobe" "https://evermeet.cx/ffmpeg/getrelease/ffprobe/zip"; then
			break
		fi
		if [[ "$i" -eq 3 ]]; then
			echo "failed to download macOS ffmpeg after 3 attempts" >&2
			exit 1
		fi
		sleep $((5 * i))
	done
	mkdir -p "$mac_dir"
	unzip -q -o "$mac_ffmpeg_zip" -d "$mac_dir"
	unzip -q -o "$mac_ffprobe_zip" -d "$mac_dir"
	find "$mac_dir" -type f \( -name ffmpeg -o -name ffprobe \) -exec cp -f {} assets/ffmpeg/darwin-arm64/ \;
	chmod +x assets/ffmpeg/darwin-arm64/ffmpeg assets/ffmpeg/darwin-arm64/ffprobe
}

#region agent log
log_debug "fetch platforms" "B" "{\"platforms\":\"$(need_platforms)\",\"version\":\"${version}\"}"
#endregion agent log

for platform in $(need_platforms); do
	case "$platform" in
		windows) fetch_windows ;;
		linux) fetch_linux ;;
		macos) fetch_macos ;;
		*) echo "unknown platform: ${platform}" >&2; exit 1 ;;
	esac
done

echo "ffmpeg ${version} fetched into assets/ffmpeg/"
