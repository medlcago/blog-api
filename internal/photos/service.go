package photos

import (
	"blog-api/internal/database"
	"blog-api/internal/models"
	"blog-api/internal/storage"
	"blog-api/pkg/errors"
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type IPhotoService interface {
	UploadAvatar(ctx context.Context, userID uint, file *multipart.FileHeader) (*UploadAvatarResponse, error)
}
type PhotoService struct {
	db        *database.DB
	minio     *storage.MinioClient
	processor *Processor
}

func NewPhotoService(db *database.DB, minio *storage.MinioClient, processor *Processor) IPhotoService {
	return &PhotoService{
		db:        db,
		minio:     minio,
		processor: processor,
	}
}

func (s *PhotoService) UploadAvatar(ctx context.Context, userID uint, file *multipart.FileHeader) (*UploadAvatarResponse, error) {
	if err := s.processor.Validate(file); err != nil {
		return nil, errors.New(400, err.Error())
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d_%d%s", uuid.New(), time.Now().UTC().UnixNano(), userID, ext)
	url := fmt.Sprintf("http://%s/%s/%s", s.minio.Client.EndpointURL().Host, s.minio.Bucket, filename)

	db := s.db.Get().WithContext(ctx)
	_, err = s.minio.Client.PutObject(ctx, s.minio.Bucket, filename, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return nil, err
	}

	if err = db.Model(&models.User{}).Where("id = ?", userID).Update("avatar", url).Error; err != nil {
		return nil, err
	}

	return &UploadAvatarResponse{
		URL: url,
	}, nil
}
