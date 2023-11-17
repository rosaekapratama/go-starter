package storage

import (
	"cloud.google.com/go/storage"
	"context"
)

type IClient interface {
	Upload(ctx context.Context, bucketName string, filePath string, fileType *string, src []byte) (written int, err error)
	Download(ctx context.Context, bucketName string, path string) (obj *storage.ObjectHandle, data []byte, err error)
}

type ClientImpl struct {
	storageClient *storage.Client
}
