package file

import (
	"archive/tar"
	"archive/zip"
	"context"
	"fmt"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/otel"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	spanZip = "common.file.Zip"
	spanTar = "common.file.Tar"
)

const defaultUnixPermissionBits = 0755

func GetParentDirPath(filePath string) string {
	s := fmt.Sprintf("%c", os.PathSeparator)
	i := strings.LastIndex(filePath, s)
	if i < 0 {
		return str.Empty
	}
	return filePath[:i+1]
}

func MkParentDir(ctx context.Context, filePath string) error {
	dir := GetParentDirPath(filePath)
	if dir != str.Empty {
		err := os.MkdirAll(dir, defaultUnixPermissionBits)
		if err != nil {
			log.Errorf(ctx, err, "Create parent directory '%s' failed", dir)
			return err
		}
		return nil
	}
	return nil
}

func CreateFileInTempDir(ctx context.Context, fileName string) (*os.File, error) {
	fileName = filepath.FromSlash(fileName)
	path := filepath.Join(os.TempDir(), fileName)
	err := MkParentDir(ctx, path)
	if err != nil {
		log.Errorf(ctx, err, "Failed to create file in temporary directory, file=%s", path)
		return nil, err
	}
	return os.Create(path)
}

func Zip(ctx context.Context, archiveName string, files []*os.File) (*os.File, error) {
	ctx, span := otel.Trace(ctx, spanZip)
	defer span.End()

	archiveFile, err := CreateFileInTempDir(ctx, archiveName)
	if err != nil {
		log.Errorf(ctx, err, "Failed to create archive file")
		return nil, err
	}
	zipWriter := zip.NewWriter(archiveFile)
	defer func(zipWriter *zip.Writer) {
		err := zipWriter.Close()
		if err != nil {
			log.Error(ctx, err)
		}
	}(zipWriter)

	for _, file := range files {
		writer, err := zipWriter.Create(filepath.Base(file.Name()))
		if err != nil {
			log.Errorf(ctx, err, "Failed to create file in zip archive, fileName=%s, archiveName=%s", file.Name(), archiveName)
			return nil, err
		}

		log.Tracef(ctx, "Writing file to zip archive, fileName=%s, archiveName=%s", file.Name(), archiveName)
		_, err = io.Copy(writer, file)
		if err != nil {
			log.Errorf(ctx, err, "Failed to write file to zip archive, fileName=%s, archiveName=%s", file.Name(), archiveName)
			return nil, err
		}
	}

	return archiveFile, nil
}

func Tar(ctx context.Context, archiveName string, files []*os.File) (*os.File, error) {
	ctx, span := otel.Trace(ctx, spanTar)
	defer span.End()

	archiveFile, err := CreateFileInTempDir(ctx, archiveName)
	if err != nil {
		log.Errorf(ctx, err, "Failed to create archive file")
		return nil, err
	}
	tarWriter := tar.NewWriter(archiveFile)
	defer func(tarWriter *tar.Writer) {
		err := tarWriter.Close()
		if err != nil {
			log.Error(ctx, err)
		}
	}(tarWriter)

	for _, file := range files {
		fileInfo, err := file.Stat()
		if err != nil {
			log.Errorf(ctx, err, "Failed to get file info, fileName=%s, archiveName=%s", file.Name(), archiveName)
			return nil, err
		}

		header, err := tar.FileInfoHeader(fileInfo, filepath.Base(file.Name()))
		if err != nil {
			log.Errorf(ctx, err, "Failed to create tar header, fileName=%s, archiveName=%s", file.Name(), archiveName)
			return nil, err
		}

		err = tarWriter.WriteHeader(header)
		if err != nil {
			log.Errorf(ctx, err, "Failed to write tar header, fileName=%s, archiveName=%s", file.Name(), archiveName)
			return nil, err
		}

		log.Tracef(ctx, "Writing file to tar archive, fileName=%s, archiveName=%s", file.Name(), archiveName)
		_, err = io.Copy(tarWriter, file)
		if err != nil {
			log.Errorf(ctx, err, "Failed to write file to tar archive, fileName=%s, archiveName=%s", file.Name(), archiveName)
			return nil, err
		}
	}

	return archiveFile, nil
}
