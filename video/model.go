package video

import (
	"fmt"
	"net/url"
	"strings"
)

type Format int

const (
	FormatMP4 Format = iota + 1
	FormatMP3
)

func FormatFromString(s string) (Format, error) {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "1", "mp4", "video":
		return FormatMP4, nil
	case "2", "mp3", "audio":
		return FormatMP3, nil
	default:
		return 0, fmt.Errorf("formato invalido: %q (use 1=mp4 ou 2=mp3)", s)
	}
}

func (f Format) String() string {
	switch f {
	case FormatMP4:
		return "mp4"
	case FormatMP3:
		return "mp3"
	default:
		return "unknown"
	}
}

func ParseURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("url vazia")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("url invalida: %w", err)
	}
	if u.Scheme == "" {
		u, err = url.Parse("https://" + raw)
		if err != nil {
			return "", fmt.Errorf("url invalida: %w", err)
		}
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("url deve usar http ou https")
	}
	host := strings.ToLower(strings.TrimSuffix(u.Hostname(), "."))
	if !isYouTubeHost(host) {
		return "", fmt.Errorf("apenas URLs do YouTube sao suportadas")
	}
	u.Fragment = ""
	return u.String(), nil
}

func isYouTubeHost(host string) bool {
	return host == "youtube.com" ||
		host == "www.youtube.com" ||
		host == "m.youtube.com" ||
		host == "music.youtube.com" ||
		host == "youtu.be" ||
		strings.HasSuffix(host, ".youtube.com")
}
