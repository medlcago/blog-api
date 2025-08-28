package photos

import (
	"blog-api/internal/database"
	"blog-api/internal/models"
	"blog-api/internal/storage"
	"blog-api/pkg/errors"
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type IPhotoService interface {
	UploadAvatar(ctx context.Context, userID uint, file *multipart.FileHeader) (*UploadAvatarResponse, error)
}
type PhotoService struct {
	db    *database.DB
	minio *storage.MinioClient
}

func NewPhotoService(db *database.DB, minio *storage.MinioClient) IPhotoService {
	return &PhotoService{
		db:    db,
		minio: minio,
	}
}

func (s *PhotoService) UploadAvatar(ctx context.Context, userID uint, file *multipart.FileHeader) (*UploadAvatarResponse, error) {
	if err := ValidateAvatar(file); err != nil {
		return nil, errors.New(400, err.Error())
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	ext := GetFileExt(file.Filename)
	filename := fmt.Sprintf("avatars/%d/%s_%d%s", userID, uuid.NewString(), time.Now().UTC().UnixNano(), ext)
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
