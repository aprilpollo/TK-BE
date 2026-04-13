package utils

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/url"
	"strings"

	"github.com/chai2010/webp"
)

func ConvertToWebP(src io.Reader, quality float32) (io.Reader, int64, error) {
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, 0, err
	}

	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, &webp.Options{Lossless: false, Quality: quality}); err != nil {
		return nil, 0, err
	}

	return &buf, int64(buf.Len()), nil
}

func ExtractObjectName(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	
	parts := strings.SplitN(strings.TrimPrefix(u.Path, "/"), "/", 2)
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid object URL: %s", rawURL)
	}
	return parts[1], nil
}
