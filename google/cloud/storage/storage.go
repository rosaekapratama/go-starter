package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/otel"
	"github.com/rosaekapratama/go-starter/response"
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

	Client = &clientImpl{storageClient: storageClient}
	log.Info(ctx, "Google cloud storage client service is initiated")
}

func (c *clientImpl) Upload(ctx context.Context, bucketName string, filePath string, fileType *string, reader io.Reader) (written int64, err error) {
	ctx, span := otel.Trace(ctx, spanUpload)
	defer span.End()

	// Check bucket exists or not
	b := c.storageClient.Bucket(bucketName)
	_, err = b.Attrs(ctx)
	if err != nil {
		log.Errorf(ctx, err, "error on storageClient.Bucket(bucketName), bucketName=%s, filePath=%s", bucketName, filePath)
		return
	}

	// Upload an object with storage.Writer.
	writer := b.Object(filePath).NewWriter(ctx)
	defer func() {
		err := writer.Close()
		if err != nil {
			log.Errorf(ctx, err, "error on writer.Close() for upload, bucketName=%s, filePath=%s", bucketName, filePath)
		}
	}()

	// Set Content-Type manually if not empty
	if fileType != nil {
		writer.ContentType = *fileType
	}

	written, err = io.Copy(writer, reader)
	if err != nil {
		log.Errorf(ctx, err, "Upload file failed, bucket=%s, filePath=%s", bucketName, filePath)
		return
	}

	return
}

func (c *clientImpl) Download(ctx context.Context, bucketName string, filePath string, writer io.Writer) (obj *storage.ObjectHandle, err error) {
	ctx, span := otel.Trace(ctx, spanDownload)
	defer span.End()

	// Check bucket exists or not
	b := c.storageClient.Bucket(bucketName)
	_, err = b.Attrs(ctx)
	if err != nil {
		log.Error(ctx, err)
		return
	}

	// Download an object with storage.Reader.
	obj = b.Object(filePath)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		log.Errorf(ctx, err, "Failed create reader for download, bucket=%s, filePath=%s", bucketName, filePath)
		return
	}

	_, err = io.Copy(writer, reader)
	if err != nil {
		log.Errorf(ctx, err, "Failed to copy all data for download, bucket=%s, filePath=%s", bucketName, filePath)
		return
	}

	return
}

func ConstructPublicUrl(bucketName string, path string) string {
	return fmt.Sprintf("%s/%s/%s", PublicUrl, bucketName, path)
}

func ConstructAuthenticatedUrl(bucketName string, path string) string {
	return fmt.Sprintf("%s/%s/%s", AuthenticatedUrl, bucketName, path)
}

// NewStreamUpload can upload data per chunk
func (c *clientImpl) NewStreamUpload(ctx context.Context, bucketName string, filePath string, fileType *string) (stream StreamUpload, err error) {
	ctx, span := otel.Trace(ctx, spanNewStreamUpload)
	defer span.End()

	// Check bucket exists or not
	b := c.storageClient.Bucket(bucketName)
	_, err = b.Attrs(ctx)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}

	// Get object handle
	objectHandle := b.Object(filePath)

	stream = &streamUploadImpl{
		streamImpl: streamImpl{
			bucketName:   bucketName,
			filePath:     filePath,
			objectHandle: objectHandle,
		},
		fileType: fileType,
	}
	return
}

// NewStreamDownload can download data per chunk
func (c *clientImpl) NewStreamDownload(ctx context.Context, bucketName string, filePath string) (stream StreamDownload, err error) {
	ctx, span := otel.Trace(ctx, spanNewStreamDownload)
	defer span.End()

	// Check bucket exists or not
	b := c.storageClient.Bucket(bucketName)
	_, err = b.Attrs(ctx)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}

	// Get object handle
	objectHandle := b.Object(filePath)

	stream = &streamDownloadImpl{
		streamImpl: streamImpl{
			bucketName:   bucketName,
			filePath:     filePath,
			objectHandle: objectHandle,
		},
	}
	return
}

func (c *clientImpl) IsExists(ctx context.Context, bucketName string, filePath string) (isExists bool, err error) {
	// Check bucket exists or not
	b := c.storageClient.Bucket(bucketName)
	_, err = b.Attrs(ctx)
	if err != nil {
		log.Errorf(ctx, err, "error on storageClient.Bucket(bucketName), bucketName=%s, filePath=%s", bucketName, filePath)
		return
	}

	// Attempt to get the object's attributes (metadata)
	_, err = b.Object(filePath).Attrs(ctx)
	if errors.Is(err, storage.ErrObjectNotExist) {
		// Object does not exist
		return false, nil
	}
	if err != nil {
		// Another error occurred
		log.Errorf(ctx, err, "error on b.Object(filePath).Attrs(ctx), bucketName=%s, filePath=%s", bucketName, filePath)
		return
	}

	// Object exists
	return true, nil
}

// GetObjectHandle return ObjectHandle
func (s *streamUploadImpl) GetObjectHandle() (obj *storage.ObjectHandle) {
	return s.objectHandle
}

func (s *streamUploadImpl) Upload(ctx context.Context, data []byte) (written int, err error) {
	// Create object handle Writer if null
	if s.writer == nil {
		s.writer = s.objectHandle.NewWriter(ctx)

		// Set Content-Type manually if not empty
		if s.fileType != nil {
			s.writer.ContentType = *s.fileType
		}
	}

	written, err = s.writer.Write(data)
	s.totalWritten += written
	return
}

func (s *streamUploadImpl) Close() (err error) {
	if s.isClosed {
		return s.errClosed
	}

	if s.writer != nil {
		s.errClosed = s.writer.Close()
	}
	if s.totalWritten == 0 {
		s.errClosed = response.FileSizeMustBeGreaterThanZero
	}

	s.isClosed = true
	return s.errClosed
}

func (c *streamUploadImpl) IsExists(ctx context.Context) (isExists bool, err error) {
	// Attempt to get the object's attributes (metadata)
	_, err = c.objectHandle.Attrs(ctx)
	if errors.Is(err, storage.ErrObjectNotExist) {
		// Object does not exist
		return false, nil
	}
	if err != nil {
		// Another error occurred
		log.Errorf(ctx, err, "error on b.Object(filePath).Attrs(ctx), bucketName=%s, filePath=%s", c.bucketName, c.filePath)
		return
	}

	// Object exists
	return true, nil
}

// GetObjectHandle return ObjectHandle
func (s *streamDownloadImpl) GetObjectHandle() (obj *storage.ObjectHandle) {
	return s.objectHandle
}

// Download parameter buffer is used to buffer read object from GCP
// data will be returned as chunk corresponding to the buffer size
// EOF error will be returned when reach the end of file
func (s *streamDownloadImpl) Download(ctx context.Context, buffer int) (data []byte, err error) {
	if s.reader == nil {
		s.reader, err = s.objectHandle.NewReader(ctx)
		if err != nil {
			log.Errorf(ctx, err, "Failed create reader for download, bucket=%s, filePath=%s", s.bucketName, s.filePath)
			return
		}
	}

	// Read and process data in chunks
	data = make([]byte, buffer)
	n, err := s.reader.Read(data)
	s.totalRead += n
	if n < buffer {
		data = data[:n]
	}
	return
}

func (s *streamDownloadImpl) Close() (err error) {
	if s.isClosed {
		return s.errClosed
	}

	if s.reader != nil {
		s.errClosed = s.reader.Close()
	}

	s.isClosed = true
	return s.errClosed
}

func (c *streamDownloadImpl) IsExists(ctx context.Context) (isExists bool, err error) {
	// Attempt to get the object's attributes (metadata)
	_, err = c.objectHandle.Attrs(ctx)
	if errors.Is(err, storage.ErrObjectNotExist) {
		// Object does not exist
		return false, nil
	}
	if err != nil {
		// Another error occurred
		log.Errorf(ctx, err, "error on c.objectHandle.Attrs(ctx), bucketName=%s, filePath=%s", c.bucketName, c.filePath)
		return
	}

	// Object exists
	return true, nil
}
