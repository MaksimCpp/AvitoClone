package miniostorage

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOStorage struct {
	client *minio.Client
	bucket string
}

func NewMinIOStorage(
	endpoint string,
	accessKey string,
	secretKey string,
	bucket string,
	useSSL bool,
) (*MinIOStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			accessKey, secretKey, "",
		),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(
			ctx,
			bucket,
			minio.MakeBucketOptions{},
		)
		if err != nil {
			return nil, err
		}
	}

	return &MinIOStorage{
		client: client,
		bucket: bucket,
	}, nil
}

func (s *MinIOStorage) Upload(
	ctx context.Context,
	file multipart.File,
	objectName string,
	contentType string,
	size int64,
	minioEndpointImage string,
) (string, error) {
	_, err := s.client.PutObject(
		ctx,
		s.bucket,
		objectName,
		file,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf(
		"http://%s/%s/%s",
		minioEndpointImage,
		s.bucket,
		objectName,
	)

	return url, nil
}

func (s *MinIOStorage) Delete(
	ctx context.Context,
	objectName string,
) error {
	return s.client.RemoveObject(
		ctx,
		s.bucket,
		objectName,
		minio.RemoveObjectOptions{},
	)
}
