package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"booptube/assets"
)

func installToolsDir() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return ""
	}
	dir := filepath.Dir(exe)

	candidates := []string{filepath.Join(dir, "tools")}
	if runtime.GOOS == "darwin" && strings.Contains(dir, ".app"+string(filepath.Separator)+"Contents"+string(filepath.Separator)+"MacOS") {
		candidates = append(candidates, filepath.Join(dir, "..", "Resources", "tools"))
	}

	for _, c := range candidates {
		clean, err := filepath.Abs(c)
		if err != nil {
			continue
		}
		if toolsDirReady(clean) {
			return clean
		}
	}
	return ""
}

func toolsDirReady(dir string) bool {
	ytdlp := filepath.Join(dir, assets.YtdlpName)
	ffmpeg := filepath.Join(dir, assets.FfmpegName)
	ffprobe := filepath.Join(dir, assets.FfprobeName)
	return isExecutableFile(ytdlp) && isExecutableFile(ffmpeg) && isExecutableFile(ffprobe)
}

func isExecutableFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}
	if runtime.GOOS == "windows" {
		return true
	}
	return info.Mode()&0o111 != 0
}
