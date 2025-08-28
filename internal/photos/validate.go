package photos

import (
	goerrors "errors"
	"fmt"
	"io"
	"mime/multipart"
	"slices"
)

type ValidationCriteria struct {
	MaxFileSize    int64
	AllowedFormats []string
	MinWidth       int
	MinHeight      int
	MaxWidth       int
	MaxHeight      int
	MinAspectRatio float64
	MaxAspectRatio float64
}

var (
	AvatarCriteria = ValidationCriteria{
		MaxFileSize:    2 << 20, // 2MB
		AllowedFormats: []string{"image/jpeg", "image/png", "image/webp"},
		MinWidth:       400,
		MinHeight:      400,
		MaxWidth:       2560,
		MaxHeight:      2560,
		MinAspectRatio: 0.25,
		MaxAspectRatio: 3.0,
	}
)

func ValidatePhoto(fileHeader *multipart.FileHeader, criteria ValidationCriteria) error {
	if fileHeader.Size > criteria.MaxFileSize {
		return fmt.Errorf("file size exceeds the limit: %d", criteria.MaxFileSize)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	contentType, err := DetectContentType(file)
	if err != nil {
		return goerrors.New("failed to detect content type")
	}

	if !slices.Contains(criteria.AllowedFormats, contentType) {
		return fmt.Errorf("invalid content type: %s", contentType)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	img, err := DecodeImageConfig(file, contentType)
	if err != nil {
		return err
	}

	if img.Width < criteria.MinWidth || img.Height < criteria.MinHeight {
		return fmt.Errorf("image dimensions too small: MinWidth: %d, MinHeight: %d", criteria.MinWidth, criteria.MinHeight)
	}

	if img.Width > criteria.MaxWidth || img.Height > criteria.MaxHeight {
		return fmt.Errorf("image dimensions too large: MaxWidth: %d, MaxHeight: %d", criteria.MaxWidth, criteria.MaxHeight)
	}

	aspectRatio := float64(img.Width) / float64(img.Height)
	if aspectRatio < criteria.MinAspectRatio || aspectRatio > criteria.MaxAspectRatio {
		return fmt.Errorf("aspect ratio must be between %.2f and %.2f", criteria.MinAspectRatio, criteria.MaxAspectRatio)
	}

	return nil
}

func ValidateAvatar(fileHeader *multipart.FileHeader) error {
	return ValidatePhoto(fileHeader, AvatarCriteria)
}
