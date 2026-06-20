//go:build darwin && !bundled

package assets

import _ "embed"

//go:embed ytdlp/darwin-arm64/yt-dlp
var Ytdlp []byte

const YtdlpName = "yt-dlp"
