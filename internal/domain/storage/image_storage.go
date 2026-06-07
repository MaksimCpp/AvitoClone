package storage

import (
	"context"
	"mime/multipart"
)

type ImageStorage interface {
	Upload(
		ctx context.Context,
		file multipart.File,
		objectName string,
		contentType string,
		size int64,
		minioEndpoint string,
	) (string, error)

	Delete(
		ctx context.Context,
		objectName string,
	) error
}