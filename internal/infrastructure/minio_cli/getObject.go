package minio_cli

import (
	"context"
	"github.com/minio/minio-go/v7"
)

func (m *minioClient) GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error) {
	return m.Client.GetObject(ctx, bucketName, objectName, opts)
}
