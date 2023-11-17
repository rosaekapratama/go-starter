package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/otel"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"io"
)

var (
	Client IClient
)

func Init(ctx context.Context, credentials *google.Credentials) {
	storageClient, err := storage.NewClient(ctx, option.WithCredentials(credentials))
	if err != nil {
		log.Fatal(ctx, err, "Failed to create google storage client")
		return
	}

	Client = &ClientImpl{storageClient: storageClient}
	log.Info(ctx, "Google cloud storage client service is initiated")
}

func (c *ClientImpl) Upload(ctx context.Context, bucketName string, filePath string, fileType *string, src []byte) (written int, err error) {
	ctx, span := otel.Trace(ctx, spanUploadFile)
	defer span.End()

	// Check bucket exists or not
	b := c.storageClient.Bucket(bucketName)
	_, err = b.Attrs(ctx)
	if err != nil {
		log.Error(ctx, err)
		return integer.Zero, err
	}

	// Upload an object with storage.Writer.
	w := b.Object(filePath).NewWriter(ctx)
	defer func(w *storage.Writer) {
		err := w.Close()
		if err != nil {
			log.Errorf(ctx, err, "Upload file fail to close writer")
		}
	}(w)

	// Set Content-Type manually if not empty
	if fileType != nil {
		w.ContentType = *fileType
	}

	written, err = w.Write(src)
	if err != nil {
		log.Errorf(ctx, err, "Upload file failed, bucket=%s, filePath=%s", bucketName, filePath)
		return integer.Zero, err
	}

	return written, nil
}

func (c *ClientImpl) Download(ctx context.Context, bucketName string, path string) (obj *storage.ObjectHandle, data []byte, err error) {
	ctx, span := otel.Trace(ctx, spanDownloadFile)
	defer span.End()

	// Check bucket exists or not
	b := c.storageClient.Bucket(bucketName)
	_, err = b.Attrs(ctx)
	if err != nil {
		log.Error(ctx, err)
		return
	}

	// Download an object with storage.Reader.
	obj = b.Object(path)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		log.Errorf(ctx, err, "Failed create reader for download, bucket=%s, filePath=%s", bucketName, path)
		return nil, nil, err
	}

	data, err = io.ReadAll(reader)
	if err != nil {
		log.Errorf(ctx, err, "Failed to read all data for download, bucket=%s, filePath=%s", bucketName, path)
		return nil, nil, err
	}

	return
}

func ConstructPublicUrl(bucketName string, path string) string {
	return fmt.Sprintf("%s/%s/%s", PublicUrl, bucketName, path)
}

func ConstructAuthenticatedUrl(bucketName string, path string) string {
	return fmt.Sprintf("%s/%s/%s", AuthenticatedUrl, bucketName, path)
}
