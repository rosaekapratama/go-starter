package utils

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	uuid2 "github.com/google/uuid"
	"github.com/rosaekapratama/go-starter/constant/env"
	"github.com/rosaekapratama/go-starter/constant/headers"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/constant/sym"
	myContext "github.com/rosaekapratama/go-starter/context"
	"github.com/rosaekapratama/go-starter/files"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/otel"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	spancloneHttpBody = "common.utils.cloneHttpBody"
)

func BoolP(b bool) *bool {
	return &b
}

func PBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func StringP(s string) *string {
	if s == str.Empty {
		return nil
	}
	return &s
}

func PString(s *string) string {
	if s == nil {
		return str.Empty
	}
	return *s
}

func IntP(i int) *int {
	return &i
}

func PInt(i *int) int {
	if i == nil {
		return integer.Zero
	}
	return *i
}

func Int32P(i int32) *int32 {
	return &i
}

func PInt32(i *int32) int32 {
	if i == nil {
		return integer.Zero
	}
	return *i
}

func Int64P(i int64) *int64 {
	return &i
}

func PInt64(i *int64) int64 {
	if i == nil {
		return integer.Zero
	}
	return *i
}

func UintP(i uint) *uint {
	return &i
}

func PUint(i *uint) uint {
	if i == nil {
		return integer.Zero
	}
	return *i
}

func Uint32P(i uint32) *uint32 {
	return &i
}

func PUint32(i *uint32) uint32 {
	if i == nil {
		return integer.Zero
	}
	return *i
}

func Uint64P(i uint64) *uint64 {
	return &i
}

func PUint64(i *uint64) uint64 {
	if i == nil {
		return integer.Zero
	}
	return *i
}

func Float32P(i float32) *float32 {
	return &i
}

func PFloat32(i *float32) float32 {
	if i == nil {
		return integer.Zero
	}
	return *i
}

func Float64P(i float64) *float64 {
	return &i
}

func PFloat64(i *float64) float64 {
	if i == nil {
		return integer.Zero
	}
	return *i
}

func GenerateSalt(ctx context.Context, size int) ([]byte, error) {
	salt := make([]byte, size)
	var _, err = rand.Read(salt)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	return salt, nil
}

func StringToSliceOfString(s string, sep string) []string {
	return strings.Split(s, sep)
}

func StringToSliceOfInt(ctx context.Context, s string, sep string) ([]int, error) {
	ints := make([]int, integer.Zero)
	for _, s := range strings.Split(s, sep) {
		i, err := strconv.Atoi(s)
		if err != nil {
			log.Errorf(ctx, err, "Failed parse '%s' to uint64", s)
			return nil, err
		}
		ints = append(ints, i)
	}
	return ints, nil
}

func StringToSliceOfInt16(s string, sep string) ([]int16, error) {
	uints := make([]int16, integer.Zero)
	for _, s := range strings.Split(s, sep) {
		if s == str.Empty {
			continue
		}
		i, err := strconv.ParseInt(s, integer.Ten, integer.I16)
		if err != nil {
			log.Errorf(context.Background(), err, "Failed parse '%s' to uint64", s)
			return nil, err
		}
		uints = append(uints, int16(i))
	}
	return uints, nil
}

func StringToSliceOfUint64(s string, sep string) ([]uint64, error) {
	uints := make([]uint64, integer.Zero)
	for _, s := range strings.Split(s, sep) {
		if s == str.Empty {
			continue
		}
		i, err := strconv.ParseUint(s, integer.Ten, integer.I64)
		if err != nil {
			log.Errorf(context.Background(), err, "Failed parse '%s' to uint64", s)
			return nil, err
		}
		uints = append(uints, i)
	}
	return uints, nil
}

func StringToSliceOfFloat64(s string, sep string) ([]float64, error) {
	uints := make([]float64, integer.Zero)
	for _, s := range strings.Split(s, sep) {
		if s == str.Empty {
			continue
		}
		i, err := strconv.ParseFloat(s, integer.I64)
		if err != nil {
			log.Errorf(context.Background(), err, "Failed parse '%s' to uint64", s)
			return nil, err
		}
		uints = append(uints, i)
	}
	return uints, nil
}

func SliceOfStringToString(strs []string, sep string) string {
	v := strings.Builder{}
	for i, s := range strs {
		v.WriteString(s)
		if i < len(strs)-1 {
			v.WriteString(sep)
		}
	}
	return v.String()
}

func SliceOfIntToString(ints []int, sep string) string {
	v := strings.Builder{}
	for idx, in := range ints {
		v.WriteString(strconv.Itoa(in))
		if idx < len(ints)-1 {
			v.WriteString(sep)
		}
	}
	return v.String()
}

func SliceOfInt16ToString(ints []int16, sep string) string {
	v := strings.Builder{}
	for i, u := range ints {
		v.WriteString(strconv.FormatInt(int64(u), integer.Ten))
		if i < len(ints)-1 {
			v.WriteString(sep)
		}
	}
	return v.String()
}

func SliceOfUint64ToString(uints []uint64, sep string) string {
	v := strings.Builder{}
	for i, u := range uints {
		v.WriteString(strconv.FormatUint(u, integer.Ten))
		if i < len(uints)-1 {
			v.WriteString(sep)
		}
	}
	return v.String()
}

func SliceOfFloat64ToString(floats []float64, sep string) string {
	v := strings.Builder{}
	for i, u := range floats {
		v.WriteString(strconv.FormatFloat(u, 'f', integer.Two, integer.I64))
		if i < len(floats)-1 {
			v.WriteString(sep)
		}
	}
	return v.String()
}

func IsRunLocally(ctx context.Context) bool {
	if localRunStr, ok := os.LookupEnv(env.EnvLocalRun); localRunStr != str.Empty && ok {
		localRun, err := strconv.ParseBool(localRunStr)
		if err != nil {
			log.Warnf(ctx, "Failed to parse %s env var '%s' to boolean, %s", env.EnvLocalRun, localRunStr, err.Error())
		} else {
			return localRun
		}
	}

	return false
}

func IsZeroValue(v interface{}) bool {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return true
		}
		val = val.Elem()
	}
	zeroValue := reflect.Zero(val.Type()).Interface()
	return reflect.DeepEqual(val.Interface(), zeroValue)
}

// CloneHttpRequest body request will be saved in temporary file, that temp file needs to be closed in defer manner
func CloneHttpRequest(req *http.Request, payloadSizeLimit int) (clonedReq *http.Request, clearFunc func(), err error) {
	ctx := req.Context()
	log.Trace(ctx, "Start clone incoming http request")

	// Clone the request
	clonedReq = new(http.Request)

	// Copy URL
	clonedReq.URL = req.URL

	// Copy method
	clonedReq.Method = req.Method

	// Copy and trim header which has length > 1KB
	// Except authorization header copy all
	clonedReq.Header = make(http.Header, len(req.Header))
	for k, v := range req.Header {
		for _, s := range v {
			if k != headers.Authorization && len(s) > payloadSizeLimit {
				clonedReq.Header.Add(k, s[:payloadSizeLimit])
			} else {
				clonedReq.Header.Add(k, s)
			}
		}
	}

	// Get content length
	contentLength := req.ContentLength
	log.Tracef(ctx, "Content length: %d", contentLength)
	log.Tracef(ctx, "Payload size limit: %d", payloadSizeLimit)

	// Process clone http body
	if req.Body != nil {
		req.Body, clonedReq.Body, clearFunc, err = cloneHttpBody(ctx, contentLength, int64(payloadSizeLimit), req.Body)
		if err != nil {
			log.Error(ctx, err, "error on cloneHttpBody(ctx, contentLength, payloadSizeLimit, req.Body)")
			return
		}
	}

	clonedReq = clonedReq.WithContext(myContext.NewContextFromTraceParent(req.Context()))
	log.Trace(ctx, "Finish clone incoming http request")
	return
}

func CloneHttpResponse(res *http.Response, payloadSizeLimit int) (clonedRes *http.Response, clearFunc func(), err error) {
	ctx := res.Request.Context()

	clonedRes = new(http.Response)

	// Copy and trim header which has length > 1KB
	// Except authorization header copy all
	clonedRes.Header = make(http.Header, len(res.Header))
	for k, v := range res.Header {
		for _, s := range v {
			if k != headers.Authorization && len(s) > payloadSizeLimit {
				clonedRes.Header.Add(k, s[:payloadSizeLimit])
			} else {
				clonedRes.Header.Add(k, s)
			}
		}
	}

	// Get content length
	contentLength := res.ContentLength

	// Process clone http body
	if res.Body != nil {
		res.Body, clonedRes.Body, clearFunc, err = cloneHttpBody(ctx, contentLength, int64(payloadSizeLimit), res.Body)
		if err != nil {
			log.Error(ctx, err, "error on cloneHttpBody(ctx, contentLength, payloadSizeLimit, req.Body)")
			return
		}
	}

	return
}

func cloneHttpBody(ctx context.Context, contentLength int64, payloadSizeLimit int64, body io.ReadCloser) (oriHttpBody io.ReadCloser, clonedHttpBody io.ReadCloser, clearFunc func(), err error) {
	ctx, span := otel.Trace(ctx, spancloneHttpBody)
	defer span.End()

	traceId := myContext.TraceIdFromContext(ctx)
	spanId := myContext.SpanIdFromContext(ctx)
	uuid := uuid2.NewString()

	// This local files is for storing body if its size > payload limit
	var oriFileBody *os.File
	var clonedFileBody *os.File
	clearFunc = func() {
		if oriFileBody != nil {
			err := oriFileBody.Close()
			if err != nil {
				log.Warn(ctx, err)
			}
			err = os.Remove(oriFileBody.Name())
			if err != nil {
				log.Warn(ctx, err)
			}
		}
		if clonedFileBody != nil {
			err := clonedFileBody.Close()
			if err != nil {
				log.Warn(ctx, err)
			}
			err = os.Remove(clonedFileBody.Name())
			if err != nil {
				log.Warn(ctx, err)
			}
		}
	}

	// If Content-Length header is < 1, it means the size is unknown,
	// then look up the size by read the body to temporary local file.
	// Or if Content-Length header is exists and > payload size limit,
	// then copy to temp local file to save memory usage.
	if contentLength < 1 || contentLength > payloadSizeLimit {

		// Create a temporary file to store original request body
		oriFileName := fmt.Sprintf("rest-o-%s-%s-%s", traceId, spanId, uuid)
		oriFileBody, err = files.CreateFileInTempDir(ctx, oriFileName)
		if err != nil {
			log.Error(ctx, err)
			return
		}

		// Store request body to temp file
		contentLength, err = io.Copy(oriFileBody, body)
		if err != nil {
			log.Error(ctx, err)
			return
		}

		// Sync the file to ensure data is written to disk
		err = oriFileBody.Sync()
		if err != nil {
			log.Error(ctx, err)
			return
		}

		// After writing to oriFileBody, rewind the file pointer back to the beginning
		_, err = oriFileBody.Seek(0, io.SeekStart)
		if err != nil {
			log.Error(ctx, err)
			return
		}

		log.Tracef(ctx, "Content length is < 1 or > %d, store original body payload to temporary local file, fileSize=%d, filePath=%s", payloadSizeLimit, contentLength, oriFileBody.Name())
	}

	// Create cloned file based on whether body size is > limit or not
	// If == 0, then set http.NoBody
	// If > limit, then use temp local file to store cloned body
	// If < limit, then use memory to store cloned body
	if contentLength == 0 {
		log.Trace(ctx, "Content length is 0, set http.NoBody")
		clonedHttpBody = http.NoBody
		oriHttpBody = http.NoBody
	} else if contentLength > payloadSizeLimit {
		// Create a temporary file to store cloned request body
		clonedFileName := fmt.Sprintf("rest-c-%s-%s-%s", traceId, spanId, uuid)
		clonedFileBody, err = files.CreateFileInTempDir(ctx, clonedFileName)
		if err != nil {
			log.Error(ctx, err)
			return
		}
		log.Tracef(ctx, "Content length is > %d, clone to temporary local file, filePath=%s", payloadSizeLimit, clonedFileBody.Name())

		// Copy trimmed body to temp file for cloned body
		trimmedBody := io.LimitReader(oriFileBody, payloadSizeLimit)
		_, err = io.Copy(clonedFileBody, trimmedBody)
		if err != nil {
			log.Error(ctx, err)
			return
		}

		// Add ellipsis
		_, err = io.WriteString(clonedFileBody, sym.Ellipsis)
		if err != nil {
			log.Error(ctx, err)
			return
		}

		// Sync the file to ensure data is written to disk
		err = clonedFileBody.Sync()
		if err != nil {
			log.Error(ctx, err)
			return
		}

		// After being read for clonedFileBody, rewind the ori file pointer back to the beginning
		_, err = oriFileBody.Seek(0, io.SeekStart)
		if err != nil {
			log.Error(ctx, err)
			return
		}

		// After writing to clonedFileBody, rewind the clone file pointer back to the beginning
		_, err = clonedFileBody.Seek(0, io.SeekStart)
		if err != nil {
			log.Error(ctx, err)
			return
		}

		// Set cloned body temp file to cloned http.Request
		clonedHttpBody = clonedFileBody

		// Restore the Body in the original http.Request
		oriHttpBody = oriFileBody
	} else if oriFileBody != nil {
		log.Tracef(ctx, "Content length is <= %d, clone to memory from temp local file", payloadSizeLimit)
		bodyBytes, errElse := io.ReadAll(oriFileBody)
		if errElse != nil {
			err = errElse
			log.Error(ctx, err)
			return
		}

		// After being read, rewind the ori file pointer back to the beginning
		_, err = oriFileBody.Seek(0, io.SeekStart)
		if err != nil {
			log.Error(ctx, err)
			return
		}

		// Set cloned body from memory to cloned http.Request
		clonedHttpBody = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Restore the Body in the original http.Request
		oriHttpBody = oriFileBody
	} else {
		log.Tracef(ctx, "Content length is <= %d, clone to memory from memory", payloadSizeLimit)
		bodyBytes, errElse := io.ReadAll(body)
		if errElse != nil {
			err = errElse
			log.Error(ctx, errElse)
			return
		}
		clonedHttpBody = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Restore the Body in the original http.Request
		oriHttpBody = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	return
}

func ChunkBytes(data []byte, chunkSize int) [][]byte {
	var chunks [][]byte

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunk := data[i:end]
		chunks = append(chunks, chunk)
	}

	return chunks
}
