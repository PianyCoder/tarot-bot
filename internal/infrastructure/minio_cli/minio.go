package minio_cli

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ObjectRepoAdapter interface {
	GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
}

type minioClient struct {
	*minio.Client
}

func NewClient(url, user, password string, port int, useSSL bool) (ObjectRepoAdapter, error) {
	endpoint := fmt.Sprintf("%s:%d", url, port)

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(user, password, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось создать структуру минио клиента: %w", err)
	}
	return &minioClient{client}, nil
}
