package drive

import (
	"context"
	"google.golang.org/api/drive/v3"
)

type IService interface {
	DownloadByFileUrl(ctx context.Context, url string) (filePath string, err error)
	DownloadByFileId(ctx context.Context, fileId string) (filePath string, err error)
	downloadFileChunk(ctx context.Context, fileId string, start int64, end int64) (chunk []byte, err error)
}

type serviceImpl struct {
	service *drive.Service
}
