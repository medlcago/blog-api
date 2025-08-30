package photos

import (
	"blog-api/internal/database"
	"blog-api/internal/errors"
	"blog-api/internal/logger"
	"blog-api/internal/models"
	"blog-api/internal/storage"
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type IPhotoService interface {
	UploadAvatar(ctx context.Context, userID uint, file *multipart.FileHeader) (*UploadAvatarResponse, error)
}
type PhotoService struct {
	db     *database.DB
	minio  *storage.MinioClient
	logger *slog.Logger
}

func NewPhotoService(db *database.DB, minio *storage.MinioClient, logger *slog.Logger) IPhotoService {
	return &PhotoService{
		db:     db,
		minio:  minio,
		logger: logger,
	}
}

func (s *PhotoService) UploadAvatar(ctx context.Context, userID uint, file *multipart.FileHeader) (*UploadAvatarResponse, error) {
	log := logger.FromCtx(ctx, s.logger).With(slog.Any("user_id", userID))

	log.Info("starting avatar upload")

	if err := ValidateAvatar(file); err != nil {
		log.Warn("avatar validation failed", logger.Err(err))
		return nil, errors.New(400, err.Error())
	}

	src, err := file.Open()
	if err != nil {
		log.Error("failed to open uploaded file", logger.Err(err))
		return nil, err
	}
	defer func() {
		if err := src.Close(); err != nil {
			log.Warn("failed to close file source", logger.Err(err))
		}
	}()

	ext := GetFileExt(file.Filename)
	filename := fmt.Sprintf("avatars/%d/%s_%d%s", userID, uuid.NewString(), time.Now().UTC().UnixNano(), ext)
	url := fmt.Sprintf("http://%s/%s/%s", s.minio.Client.EndpointURL().Host, s.minio.Bucket, filename)

	db := s.db.Get().WithContext(ctx)

	log.Info("uploading file to minio storage")
	_, err = s.minio.Client.PutObject(ctx, s.minio.Bucket, filename, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})

	if err != nil {
		log.Error("failed to upload file to minio", logger.Err(err))
		return nil, err
	}

	if err = db.Model(&models.User{}).Where("id = ?", userID).Update("avatar", url).Error; err != nil {
		log.Error("failed to update user avatar in database", logger.Err(err))

		log.Warn("attempting to rollback: delete uploaded file from minio")
		if delErr := s.minio.Client.RemoveObject(ctx, s.minio.Bucket, filename, minio.RemoveObjectOptions{}); delErr != nil {
			log.Error("failed to delete file during rollback", slog.Any("error", delErr))
		}

		return nil, err
	}

	log.Info("avatar uploaded successfully", slog.String("url", url))

	return &UploadAvatarResponse{
		URL: url,
	}, nil
}
