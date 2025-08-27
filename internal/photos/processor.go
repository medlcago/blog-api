package photos

import (
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"slices"
)

var (
	ErrInvalidFormat = errors.New("invalid file format")
	ErrTooLarge      = errors.New("file too large")
)

type Processor struct {
	MaxSize int64
	Allowed []string
}

func NewProcessor(maxSize int64, allowed []string) *Processor {
	return &Processor{MaxSize: maxSize, Allowed: allowed}
}

func (p *Processor) Validate(file *multipart.FileHeader) error {
	if file.Size > p.MaxSize {
		return ErrTooLarge
	}

	f, err := file.Open()
	defer f.Close()

	_, format, err := image.Decode(f)
	if err != nil {
		return ErrInvalidFormat
	}

	if !slices.Contains(p.Allowed, format) {
		return ErrInvalidFormat
	}

	return nil
}
