//go:build darwin

package assets

import _ "embed"

//go:embed ffmpeg/darwin-arm64/ffmpeg
var Ffmpeg []byte

//go:embed ffmpeg/darwin-arm64/ffprobe
var Ffprobe []byte

const FfmpegName = "ffmpeg"
const FfprobeName = "ffprobe"
