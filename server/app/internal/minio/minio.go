package minio

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"time"
	"uptimatic/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioUtil struct {
	Client *minio.Client
	Bucket string
}

func NewMinioUtil(ctx context.Context, cfg *config.Config) (*MinioUtil, error) {
	client, err := minio.New(cfg.StorageEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.StorageAccessKey, cfg.StorageSecretKey, ""),
		Secure: cfg.StorageUseSSL,
	})
	if err != nil {
		return nil, err
	}

	exists, err := client.BucketExists(ctx, cfg.StorageBucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = client.MakeBucket(ctx, cfg.StorageBucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &MinioUtil{Client: client, Bucket: cfg.StorageBucket}, nil
}

func (m *MinioUtil) UploadFile(ctx context.Context, file multipart.File, fileName, contentType string, size int64) error {
	_, err := m.Client.PutObject(ctx, m.Bucket, fileName, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *MinioUtil) DeleteFile(ctx context.Context, fileName string) error {
	err := m.Client.RemoveObject(ctx, m.Bucket, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (m *MinioUtil) GetPresignedURL(ctx context.Context, fileName string) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := m.Client.PresignedGetObject(ctx, m.Bucket, fileName, time.Hour, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func (m *MinioUtil) PutPresignedURL(ctx context.Context, fileName, contentType string) (string, error) {
	if contentType != "image/jpeg" &&
		contentType != "image/png" &&
		contentType != "image/webp" {
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}

	reqParams := make(url.Values)
	reqParams.Set("Content-Type", contentType)

	presignedURL, err := m.Client.PresignedPutObject(ctx, m.Bucket, fileName, time.Minute*15)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.String(), nil
}
