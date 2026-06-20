//go:build windows && !bundled

package assets

import _ "embed"

//go:embed ytdlp/windows-amd64/yt-dlp.exe
var Ytdlp []byte

const YtdlpName = "yt-dlp.exe"
