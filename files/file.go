package files

import (
	"archive/tar"
	"archive/zip"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/rosaekapratama/go-starter/constant/headers/contenttype"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/otel"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	spanZip = "common.file.Zip"
	spanTar = "common.file.Tar"
)

const defaultUnixPermissionBits = 0755

const FileExtPdf string = "pdf"
const FileExtZip string = "zip"
const FileExtGz string = "gz"
const FileExtJpg string = "jpg"
const FileExtPng string = "png"
const FileExtGif string = "gif"
const FileExtTxt string = "txt"

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

func Zip(ctx context.Context, archiveName string, filePaths []string) (zipFile string, err error) {
	ctx, span := otel.Trace(ctx, spanZip)
	defer span.End()

	archiveFile, err := CreateFileInTempDir(ctx, archiveName)
	if err != nil {
		log.Errorf(ctx, err, "Failed to create archive file, archiveName=%s", archiveName)
		return
	}
	defer func() {
		err := archiveFile.Close()
		if err != nil {
			log.Error(ctx, err, "error on archiveFile.Close(), archiveName=%s", archiveName)
		}
	}()

	zipWriter := zip.NewWriter(archiveFile)
	defer func() {
		err := zipWriter.Close()
		if err != nil {
			log.Error(ctx, err, "error on zipWriter.Close(), archiveName=%s", archiveName)
		}
	}()

	fileList := make([]*os.File, 0)
	defer func() {
		for _, file := range fileList {
			err := file.Close()
			if err != nil {
				log.Error(ctx, err, "error on file.Close(), fileName=%s, archiveName=%s", file.Name(), archiveName)
			}
		}
	}()

	for _, filePath := range filePaths {
		var file *os.File
		file, err = os.Open(filePath)
		if err != nil {
			log.Errorf(ctx, err, "Failed to open file, fileName=%s, archiveName=%s", file.Name(), archiveName)
			return
		}
		fileList = append(fileList, file)

		writer, err := zipWriter.Create(filepath.Base(file.Name()))
		if err != nil {
			log.Errorf(ctx, err, "Failed to create file in zip archive, fileName=%s, archiveName=%s", file.Name(), archiveName)
			return str.Empty, err
		}

		log.Tracef(ctx, "Writing file to zip archive, fileName=%s, archiveName=%s", file.Name(), archiveName)
		_, err = io.Copy(writer, file)
		if err != nil {
			log.Errorf(ctx, err, "Failed to write file to zip archive, fileName=%s, archiveName=%s", file.Name(), archiveName)
			return str.Empty, err
		}
	}

	return archiveFile.Name(), nil
}

func Tar(ctx context.Context, archiveName string, filePaths []string) (tarFile string, err error) {
	ctx, span := otel.Trace(ctx, spanTar)
	defer span.End()

	archiveFile, err := CreateFileInTempDir(ctx, archiveName)
	if err != nil {
		log.Errorf(ctx, err, "Failed to create archive file")
		return
	}
	defer func() {
		err := archiveFile.Close()
		if err != nil {
			log.Error(ctx, err, "error on archiveFile.Close(), archiveName=%s", archiveName)
		}
	}()

	tarWriter := tar.NewWriter(archiveFile)
	defer func() {
		err := tarWriter.Close()
		if err != nil {
			log.Error(ctx, err, "error on tarWriter.Close(), archiveName=%s", archiveName)
		}
	}()

	fileList := make([]*os.File, 0)
	defer func() {
		for _, file := range fileList {
			err := file.Close()
			if err != nil {
				log.Error(ctx, err, "error on file.Close(), fileName=%s, archiveName=%s", file.Name(), archiveName)
			}
		}
	}()

	for _, filePath := range filePaths {
		var file *os.File
		file, err = os.Open(filePath)
		if err != nil {
			log.Errorf(ctx, err, "Failed to open file, fileName=%s, archiveName=%s", file.Name(), archiveName)
			return
		}
		fileList = append(fileList, file)

		fileInfo, err := file.Stat()
		if err != nil {
			log.Errorf(ctx, err, "Failed to get file info, fileName=%s, archiveName=%s", file.Name(), archiveName)
			return str.Empty, err
		}

		header, err := tar.FileInfoHeader(fileInfo, filepath.Base(file.Name()))
		if err != nil {
			log.Errorf(ctx, err, "Failed to create tar header, fileName=%s, archiveName=%s", file.Name(), archiveName)
			return str.Empty, err
		}

		err = tarWriter.WriteHeader(header)
		if err != nil {
			log.Errorf(ctx, err, "Failed to write tar header, fileName=%s, archiveName=%s", file.Name(), archiveName)
			return str.Empty, err
		}

		log.Tracef(ctx, "Writing file to tar archive, fileName=%s, archiveName=%s", file.Name(), archiveName)
		_, err = io.Copy(tarWriter, file)
		if err != nil {
			log.Errorf(ctx, err, "Failed to write file to tar archive, fileName=%s, archiveName=%s", file.Name(), archiveName)
			return str.Empty, err
		}
	}

	return archiveFile.Name(), nil
}

func EncodeBase64File(ctx context.Context, inputFilePath string, outputFilePath string, chunkSize int) (err error) {
	// Open the input file
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		log.Error(ctx, err)
		return
	}
	defer func(inputFile *os.File) {
		err := inputFile.Close()
		if err != nil {
			log.Error(ctx, err, "error on inputFile.Close() for EncodeBase64File()")
		}
	}(inputFile)

	// Create if not exists or open if exists the output file
	outputFile, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Error(ctx, err)
		return
	}
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {
			log.Error(ctx, err, "error on outputFile.Close() for EncodeBase64File()")
		}
	}(outputFile)

	// Create a new base64 encoder
	encoder := base64.NewEncoder(base64.StdEncoding, outputFile)
	defer func(encoder io.WriteCloser) {
		err := encoder.Close()
		if err != nil {
			log.Error(ctx, err, "error on encoder.Close() for EncodeBase64File()")
		}
	}(encoder)

	// Buffer to hold chunks of data
	buffer := make([]byte, chunkSize)

	for {
		// Read chunk from the input file
		n, errFor := inputFile.Read(buffer)
		if errFor != nil && errFor != io.EOF {
			log.Error(ctx, err, "error on inputFile.Read(buffer) for EncodeBase64File()")
			err = errFor
			return
		}

		// Write the encoded data to the output file
		if n > 0 {
			if _, errFor := encoder.Write(buffer[:n]); err != nil {
				log.Error(ctx, errFor, "error on encoder.Write(buffer[:n]) for EncodeBase64File()")
				err = errFor
				return
			}
		}

		// Break if EOF is reached
		if errFor == io.EOF {
			break
		}
	}

	return
}

func DecodeBase64File(ctx context.Context, inputFilePath string, outputFilePath string, chunkSize int) (err error) {
	// Open the input file
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		log.Error(ctx, err)
		return
	}
	defer func(inputFile *os.File) {
		err := inputFile.Close()
		if err != nil {
			log.Error(ctx, err, "error on inputFile.Close() for DecodeBase64File()")
		}
	}(inputFile)

	// Create if not exists or open if exists the output file
	outputFile, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Error(ctx, err)
		return
	}
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {
			log.Error(ctx, err, "error on outputFile.Close() for DecodeBase64File()")
		}
	}(outputFile)

	// Create a new base64 decoder
	decoder := base64.NewDecoder(base64.StdEncoding, inputFile)

	// Buffer to hold chunks of data
	buffer := make([]byte, chunkSize)

	for {
		// Read chunk from the decoder
		n, errFor := decoder.Read(buffer)
		if errFor != nil && errFor != io.EOF {
			log.Error(ctx, errFor, "error on decoder.Read(buffer) for DecodeBase64File()")
			err = errFor
			return
		}

		// Write the decoded data to the output file
		if n > 0 {
			if _, errFor := outputFile.Write(buffer[:n]); err != nil {
				log.Error(ctx, errFor, "error on outputFile.Write(buffer[:n]) for EncodeBase64File()")
				err = errFor
				return
			}
		}

		// Break if EOF is reached
		if errFor == io.EOF {
			break
		}
	}

	return nil
}

func DetectHttpContentType(ctx context.Context, filePath string) (contentType string) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Errorf(ctx, err, "error on os.Open(filename), fileName=%s", filePath)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Errorf(ctx, err, "error on file.Close(), fileName=%s", filePath)
		}
	}(file)

	// Read the first 512 bytes (HTTP standard for content sniffing)
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		log.Errorf(ctx, err, "error on file.Read(buf), fileName=%s", filePath)
		return
	}

	// Detect the content type
	contentType = http.DetectContentType(buf)
	return
}

func DetectExtension(ctx context.Context, filePath string) (ext string, found bool) {
	// Detect the content type
	contentType := DetectHttpContentType(ctx, filePath)

	// Map of content types to file extensions
	extMap := map[string]string{
		contenttype.ApplicationPdf:   FileExtPdf,
		contenttype.ApplicationZip:   FileExtZip,
		contenttype.ApplicationXGzip: FileExtGz,
		contenttype.ImageJpeg:        FileExtJpg,
		contenttype.ImagePng:         FileExtPng,
		contenttype.ImageGif:         FileExtGif,
		contenttype.TextPlain:        FileExtTxt,
		// Add more mappings as needed
	}

	// Return the extension if found, otherwise default to .bin
	ext, found = extMap[contentType]
	return
}
