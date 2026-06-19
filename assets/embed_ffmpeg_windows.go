//go:build windows

package assets

import _ "embed"

//go:embed ffmpeg/windows-amd64/ffmpeg.exe
var Ffmpeg []byte

//go:embed ffmpeg/windows-amd64/ffprobe.exe
var Ffprobe []byte

const FfmpegName = "ffmpeg.exe"
const FfprobeName = "ffprobe.exe"
