//go:build linux && !bundled

package assets

import _ "embed"

//go:embed ytdlp/linux-amd64/yt-dlp
var Ytdlp []byte

const YtdlpName = "yt-dlp"
