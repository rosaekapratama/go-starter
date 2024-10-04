package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
)

type IClient interface {
	Upload(ctx context.Context, bucketName string, filePath string, fileType *string, reader io.Reader) (written int64, err error)
	Download(ctx context.Context, bucketName string, filePath string, writer io.Writer) (obj *storage.ObjectHandle, err error)
	NewStreamUpload(ctx context.Context, bucketName string, filePath string, fileType *string) (stream StreamUpload, err error)
	NewStreamDownload(ctx context.Context, bucketName string, filePath string) (stream StreamDownload, err error)
	IsExists(ctx context.Context, bucketName string, filePath string) (isExists bool, err error)
}

type clientImpl struct {
	storageClient *storage.Client
}

type stream interface {
	GetObjectHandle() (obj *storage.ObjectHandle)
	IsExists(ctx context.Context) (isExists bool, err error)
}

type StreamUpload interface {
	stream
	Upload(ctx context.Context, data []byte) (written int, err error)
	Close() (err error)
}

type StreamDownload interface {
	stream
	Download(ctx context.Context, buffer int) (data []byte, err error)
	Close() (err error)
}

type streamImpl struct {
	bucketName   string
	filePath     string
	objectHandle *storage.ObjectHandle
}

type streamUploadImpl struct {
	streamImpl
	fileType     *string
	writer       *storage.Writer
	totalWritten int
	errClosed    error
	isClosed     bool
}

type streamDownloadImpl struct {
	streamImpl
	reader    *storage.Reader
	totalRead int
	errClosed error
	isClosed  bool
}
