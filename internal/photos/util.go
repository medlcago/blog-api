package photos

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"

	"golang.org/x/image/webp"
)

func DecodeImageConfig(r io.Reader, contentType string) (image.Config, error) {
	if contentType == "image/webp" {
		return webp.DecodeConfig(r)
	}
	img, _, err := image.DecodeConfig(r)
	return img, err
}

func DetectContentType(r io.Reader) (string, error) {
	buff := make([]byte, 512)
	_, err := r.Read(buff)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buff)
	return contentType, nil
}
