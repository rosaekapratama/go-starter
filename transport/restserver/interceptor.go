package restserver

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-http-utils/headers"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/constant/sym"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/page"
	"github.com/rosaekapratama/go-starter/response"
)

const (
	noWritten                  = -1
	dataVarName                = "${data}"
	dataVarNameWithDoubleQuote = "\"${data}\""
	nullString                 = "null"
	removedJsonFieldData       = ",\"data\":\"${data}\""
)

type BaseResponse struct {
	Response   *Response          `json:"response"`
	Pagination *page.PageResponse `json:"pagination,omitempty"`
	Data       string             `json:"data,omitempty"`
}

type Response struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type WriterInterceptor struct {
	http.ResponseWriter
	size   int
	status int

	ctx       context.Context
	response  response.IResponse
	page      *page.PageResponse
	body      []byte
	IsWritten bool
	IsRaw     bool
	a         []any
}

func (w *WriterInterceptor) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *WriterInterceptor) reset(writer http.ResponseWriter) {
	w.ResponseWriter = writer
}

// WriteHeader wrapper func of http.ResponseWriter WriteHeader
func (w *WriterInterceptor) WriteHeader(statusCode int) {
	if w.response != nil && w.response.HttpStatusCode() > integer.Zero {
		w.status = w.response.HttpStatusCode()
	} else {
		w.status = statusCode
	}
	w.ResponseWriter.WriteHeader(w.status)
}

func (w *WriterInterceptor) WriteHeaderNow() {
	if w.response != nil && w.response.HttpStatusCode() > integer.Zero {
		w.status = w.response.HttpStatusCode()
	} else {
		w.status = response.Success.HttpStatusCode()
	}
	w.ResponseWriter.WriteHeader(w.status)
}

// Write wrapper func of http.ResponseWriter Write
func (w *WriterInterceptor) Write(b []byte) (int, error) {
	// Only process if content type is application/json or raw is true
	contentType := strings.ToLower(strings.Split(w.ResponseWriter.Header().Get(headers.ContentType), sym.SemiColon)[0])
	if contentType != contentTypeApplicationJson || w.IsRaw {
		w.IsWritten = true
		return w.ResponseWriter.Write(b)
	}

	if b == nil {
		b = make([]byte, integer.Zero)
	}

	r := &BaseResponse{
		Response: &Response{},
	}

	// Keep real length to be sent later
	realLen := len(b)

	// Set response to base response struct, set with response.UnknownResponse and return error if interceptor has nil response
	err := r.setResponse(w)
	if err == nil {
		// Set pagination to base response struct
		r.setPagination(w)

		// Set data with its variable name which will be replaced later with real json value
		if realLen > integer.Zero {
			r.Data = dataVarName
		}
	}

	rb, err := r.marshal(w.ctx)
	if err != nil {
		log.Error(w.ctx, err)
		w.ResponseWriter.WriteHeader(response.GeneralError.HttpStatusCode())
	} else if realLen > integer.Zero {
		// Get slice of data with optional leading whitespace removed.
		// See RFC 7159, Section 2 for the definition of JSON whitespace.
		b = bytes.TrimLeft(b, " \t\r\n")
		baseStr := string(rb)
		dataStr := string(b)
		dataStr = strings.ReplaceAll(dataStr, "\\\\r", "\\r")
		dataStr = strings.ReplaceAll(dataStr, "\\\\n", "\\n")

		var finalStr string
		if b[0] == sym.OpenSquareBracket[integer.Zero] || b[0] == sym.OpenCurlyBracket[integer.Zero] {
			// Bytes are json array or object
			finalStr = strings.ReplaceAll(baseStr, dataVarNameWithDoubleQuote, dataStr)
		} else if dataStr == nullString {
			// Remove field data from base response
			finalStr = strings.ReplaceAll(baseStr, removedJsonFieldData, str.Empty)
		} else {
			// Bytes are neither array nor object, set as a string
			finalStr = strings.ReplaceAll(baseStr, dataVarName, dataStr)
		}

		rb = []byte(finalStr)
		w.Header().Set(headers.ContentLength, strconv.Itoa(len(rb)))
	}

	w.IsWritten = true
	w.size = realLen
	w.status = w.response.HttpStatusCode()
	w.body = rb
	w.ResponseWriter.WriteHeader(w.response.HttpStatusCode())
	_, err = w.ResponseWriter.Write(rb)
	w.Header().Set(headers.ContentLength, strconv.Itoa(realLen))
	return realLen, err
}

func (w *WriterInterceptor) WriteString(s string) (n int, err error) {
	w.WriteHeaderNow()
	n, err = io.WriteString(w.ResponseWriter, s)
	w.size += n
	w.IsWritten = true
	return
}

func (w *WriterInterceptor) Status() int {
	return w.status
}

func (w *WriterInterceptor) Size() int {
	return w.size
}

func (w *WriterInterceptor) Written() bool {
	return w.size != noWritten
}

// Hijack implements the http.Hijacker interface.
func (w *WriterInterceptor) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.size < 0 {
		w.size = 0
	}
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

// CloseNotify implements the http.CloseNotifier interface.
func (w *WriterInterceptor) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Flush implements the http.Flusher interface.
func (w *WriterInterceptor) Flush() {
	w.WriteHeaderNow()
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *WriterInterceptor) Pusher() (pusher http.Pusher) {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}

func NewWriterInterceptor(ctx context.Context, w http.ResponseWriter) *WriterInterceptor {
	return &WriterInterceptor{w, integer.NOne, integer.Zero, ctx, nil, nil, nil, false, false, nil}
}

func (r *BaseResponse) setResponse(i *WriterInterceptor) error {
	if r.Response == nil {
		r.Response = &Response{}
	}

	var err error
	if i.response == nil {
		i.response = response.UnknownResponse
		err = response.UnknownResponse
	}

	r.Response.Code = i.response.Code()
	if len(i.a) > integer.Zero {
		r.Response.Description = fmt.Sprintf(i.response.Description(), i.a...)
	} else {
		r.Response.Description = i.response.Description()
	}
	return err
}

func (r *BaseResponse) setPagination(i *WriterInterceptor) {
	if i.page != nil {
		r.Pagination = i.page
	}
}

func (r *BaseResponse) marshal(ctx context.Context) ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		log.Error(ctx, err, "Post modify marshal response failed")
		r.Response.Code = response.GeneralError.Code()
		r.Response.Description = response.GeneralError.Description()
		r.Pagination = nil
		r.Data = str.Empty

		// It should not be error, just in case
		b, err = json.Marshal(&r)
		if err != nil {
			log.Error(ctx, err, "Marshal general error failed")
			return nil, err
		}
		return b, nil
	}
	return b, nil
}

func castInterceptor(w http.ResponseWriter) *WriterInterceptor {
	var i *WriterInterceptor
	switch w.(type) {
	case gin.ResponseWriter:
		i = w.(gin.ResponseWriter).(*WriterInterceptor)
	case *WriterInterceptor:
		i = w.(*WriterInterceptor)
	}
	return i
}
