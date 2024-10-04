package drive

import (
	"context"
	"fmt"
	"github.com/inhies/go-bytesize"
	"github.com/rosaekapratama/go-starter/constant/headers"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/constant/sym"
	"github.com/rosaekapratama/go-starter/files"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/otel"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"net/http"
	urlLib "net/url"
	"regexp"
)

const (
	spangetFileId = "common.google.drive.getFileId"
	spanDownload  = "common.google.drive.Download"

	hostGoogleDrive = "drive.google.com"
	hostGoogleDoc   = "docs.google.com"
	chunkSize       = 10 * int64(bytesize.MB)
)

var (
	Service IService
)

func Init(ctx context.Context, credentials *google.Credentials) {
	// Define the scopes it needs
	scopes := []string{
		drive.DriveReadonlyScope,
	}

	service, err := drive.NewService(ctx, option.WithCredentials(credentials), option.WithScopes(scopes...))
	if err != nil {
		log.Fatal(ctx, err, "Failed to create google drive service")
	}

	Service = &serviceImpl{
		service: service,
	}
}

func (m *serviceImpl) DownloadByFileUrl(ctx context.Context, url string) (filePath string, err error) {
	ctx, span := otel.Trace(ctx, spanDownload)
	defer span.End()

	fileId, err := getFileId(ctx, url)
	if err != nil {
		log.Error(ctx, err, "error on getFileId()")
		return
	}

	return m.DownloadByFileId(ctx, fileId)
}

func (m *serviceImpl) DownloadByFileId(ctx context.Context, fileId string) (filePath string, err error) {
	ctx, span := otel.Trace(ctx, spanDownload)
	defer span.End()

	// create temp file to store downloaded bytes
	file, err := files.CreateFileInTempDir(ctx, fmt.Sprintf("google.drive.%s", fileId))
	if err != nil {
		log.Error(ctx, err, "error on files.CreateFileInTempDir(), fileId=%s", fileId)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Error(ctx, err)
		}
	}()

	var start, end int64
	for {
		end = start + chunkSize - 1
		chunk, errFor := m.downloadFileChunk(ctx, fileId, start, end)
		if errFor != nil {
			err = errFor
			log.Error(ctx, errFor, "error on m.downloadFileChunk(ctx, start, end)")
			return
		}

		// if chunk == 0 means EOF
		if len(chunk) == 0 {
			break
		}

		// write to temp file
		_, errFor = file.Write(chunk)
		if errFor != nil {
			err = errFor
			log.Error(ctx, errFor, "error on file.Write(chunk)")
			return
		}

		// set next offset
		start += chunkSize
	}

	filePath = file.Name()
	return
}

func (m *serviceImpl) downloadFileChunk(ctx context.Context, fileId string, start int64, end int64) (chunk []byte, err error) {
	// prepare download request
	filesGetCall := m.service.Files.Get(fileId)
	filesGetCall.SupportsAllDrives(true)
	filesGetCall.Header().Add(headers.Range, fmt.Sprintf("bytes=%d-%d", start, end))

	// execute download
	res, err := filesGetCall.Download()
	if err != nil {
		log.Errorf(ctx, err, "error on filesGetCall.Download(), fileId=%s", fileId)
		return
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Error(ctx, err)
		}
	}()

	if res.StatusCode != http.StatusPartialContent {
		return nil, fmt.Errorf("expected status 206 Partial Content, got %d", res.StatusCode)
	}

	chunk, err = io.ReadAll(res.Body)
	return
}

// getFileId extracts the file ID from a Google Drive or Google Doc URL.
func getFileId(ctx context.Context, url string) (fileId string, err error) {
	ctx, span := otel.Trace(ctx, spangetFileId)
	defer span.End()

	parsedURL, err := urlLib.Parse(url)
	if err != nil {
		log.Errorf(ctx, err, "error on url.parse(), url=%s", url)
		return str.Empty, err
	}

	var re *regexp.Regexp
	switch parsedURL.Host {
	case hostGoogleDrive:
		// Google Drive URL patterns
		re = regexp.MustCompile(`^/file/d/([^/]+)|^/open\?id=([^&]+)|^/drive/folders/([^/]+)`)
	case hostGoogleDoc:
		// Google Docs URL pattern
		re = regexp.MustCompile(`^/document/d/([^/]+)`)
	default:
		err = fmt.Errorf("unsupported host: %s", parsedURL.Host)
		log.Error(ctx, err)
		return
	}

	// Find the first matching group
	matches := re.FindStringSubmatch(parsedURL.Path + sym.QuestionMark + parsedURL.RawQuery)
	if len(matches) > 1 {
		for _, match := range matches[1:] {
			if match != str.Empty {
				fileId = match
				return
			}
		}
	}

	err = fmt.Errorf("invalid URL: %v", err)
	log.Error(ctx, err)
	return
}
