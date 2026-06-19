//go:build linux

package assets

import _ "embed"

//go:embed ffmpeg/linux-amd64/ffmpeg
var Ffmpeg []byte

//go:embed ffmpeg/linux-amd64/ffprobe
var Ffprobe []byte

const FfmpegName = "ffmpeg"
const FfprobeName = "ffprobe"
