package utils

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

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
