// Code generated by mockery v2.46.3. DO NOT EDIT.

package storage

import (
	context "context"

	cloudstorage "github.com/rosaekapratama/go-starter/google/cloud/storage"

	io "io"

	mock "github.com/stretchr/testify/mock"

	storage "cloud.google.com/go/storage"
)

// MockIClient is an autogenerated mock type for the IClient type
type MockIClient struct {
	mock.Mock
}

type MockIClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIClient) EXPECT() *MockIClient_Expecter {
	return &MockIClient_Expecter{mock: &_m.Mock}
}

// Download provides a mock function with given fields: ctx, bucketName, filePath, writer
func (_m *MockIClient) Download(ctx context.Context, bucketName string, filePath string, writer io.Writer) (*storage.ObjectHandle, error) {
	ret := _m.Called(ctx, bucketName, filePath, writer)

	if len(ret) == 0 {
		panic("no return value specified for Download")
	}

	var r0 *storage.ObjectHandle
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, io.Writer) (*storage.ObjectHandle, error)); ok {
		return rf(ctx, bucketName, filePath, writer)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, io.Writer) *storage.ObjectHandle); ok {
		r0 = rf(ctx, bucketName, filePath, writer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.ObjectHandle)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, io.Writer) error); ok {
		r1 = rf(ctx, bucketName, filePath, writer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIClient_Download_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Download'
type MockIClient_Download_Call struct {
	*mock.Call
}

// Download is a helper method to define mock.On call
//   - ctx context.Context
//   - bucketName string
//   - filePath string
//   - writer io.Writer
func (_e *MockIClient_Expecter) Download(ctx interface{}, bucketName interface{}, filePath interface{}, writer interface{}) *MockIClient_Download_Call {
	return &MockIClient_Download_Call{Call: _e.mock.On("Download", ctx, bucketName, filePath, writer)}
}

func (_c *MockIClient_Download_Call) Run(run func(ctx context.Context, bucketName string, filePath string, writer io.Writer)) *MockIClient_Download_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(io.Writer))
	})
	return _c
}

func (_c *MockIClient_Download_Call) Return(obj *storage.ObjectHandle, err error) *MockIClient_Download_Call {
	_c.Call.Return(obj, err)
	return _c
}

func (_c *MockIClient_Download_Call) RunAndReturn(run func(context.Context, string, string, io.Writer) (*storage.ObjectHandle, error)) *MockIClient_Download_Call {
	_c.Call.Return(run)
	return _c
}

// IsExists provides a mock function with given fields: ctx, bucketName, filePath
func (_m *MockIClient) IsExists(ctx context.Context, bucketName string, filePath string) (bool, error) {
	ret := _m.Called(ctx, bucketName, filePath)

	if len(ret) == 0 {
		panic("no return value specified for IsExists")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (bool, error)); ok {
		return rf(ctx, bucketName, filePath)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctx, bucketName, filePath)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, bucketName, filePath)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIClient_IsExists_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsExists'
type MockIClient_IsExists_Call struct {
	*mock.Call
}

// IsExists is a helper method to define mock.On call
//   - ctx context.Context
//   - bucketName string
//   - filePath string
func (_e *MockIClient_Expecter) IsExists(ctx interface{}, bucketName interface{}, filePath interface{}) *MockIClient_IsExists_Call {
	return &MockIClient_IsExists_Call{Call: _e.mock.On("IsExists", ctx, bucketName, filePath)}
}

func (_c *MockIClient_IsExists_Call) Run(run func(ctx context.Context, bucketName string, filePath string)) *MockIClient_IsExists_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockIClient_IsExists_Call) Return(isExists bool, err error) *MockIClient_IsExists_Call {
	_c.Call.Return(isExists, err)
	return _c
}

func (_c *MockIClient_IsExists_Call) RunAndReturn(run func(context.Context, string, string) (bool, error)) *MockIClient_IsExists_Call {
	_c.Call.Return(run)
	return _c
}

// NewStreamDownload provides a mock function with given fields: ctx, bucketName, filePath
func (_m *MockIClient) NewStreamDownload(ctx context.Context, bucketName string, filePath string) (cloudstorage.StreamDownload, error) {
	ret := _m.Called(ctx, bucketName, filePath)

	if len(ret) == 0 {
		panic("no return value specified for NewStreamDownload")
	}

	var r0 cloudstorage.StreamDownload
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (cloudstorage.StreamDownload, error)); ok {
		return rf(ctx, bucketName, filePath)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) cloudstorage.StreamDownload); ok {
		r0 = rf(ctx, bucketName, filePath)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cloudstorage.StreamDownload)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, bucketName, filePath)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIClient_NewStreamDownload_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NewStreamDownload'
type MockIClient_NewStreamDownload_Call struct {
	*mock.Call
}

// NewStreamDownload is a helper method to define mock.On call
//   - ctx context.Context
//   - bucketName string
//   - filePath string
func (_e *MockIClient_Expecter) NewStreamDownload(ctx interface{}, bucketName interface{}, filePath interface{}) *MockIClient_NewStreamDownload_Call {
	return &MockIClient_NewStreamDownload_Call{Call: _e.mock.On("NewStreamDownload", ctx, bucketName, filePath)}
}

func (_c *MockIClient_NewStreamDownload_Call) Run(run func(ctx context.Context, bucketName string, filePath string)) *MockIClient_NewStreamDownload_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockIClient_NewStreamDownload_Call) Return(stream cloudstorage.StreamDownload, err error) *MockIClient_NewStreamDownload_Call {
	_c.Call.Return(stream, err)
	return _c
}

func (_c *MockIClient_NewStreamDownload_Call) RunAndReturn(run func(context.Context, string, string) (cloudstorage.StreamDownload, error)) *MockIClient_NewStreamDownload_Call {
	_c.Call.Return(run)
	return _c
}

// NewStreamUpload provides a mock function with given fields: ctx, bucketName, filePath, fileType
func (_m *MockIClient) NewStreamUpload(ctx context.Context, bucketName string, filePath string, fileType *string) (cloudstorage.StreamUpload, error) {
	ret := _m.Called(ctx, bucketName, filePath, fileType)

	if len(ret) == 0 {
		panic("no return value specified for NewStreamUpload")
	}

	var r0 cloudstorage.StreamUpload
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *string) (cloudstorage.StreamUpload, error)); ok {
		return rf(ctx, bucketName, filePath, fileType)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *string) cloudstorage.StreamUpload); ok {
		r0 = rf(ctx, bucketName, filePath, fileType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cloudstorage.StreamUpload)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, *string) error); ok {
		r1 = rf(ctx, bucketName, filePath, fileType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIClient_NewStreamUpload_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NewStreamUpload'
type MockIClient_NewStreamUpload_Call struct {
	*mock.Call
}

// NewStreamUpload is a helper method to define mock.On call
//   - ctx context.Context
//   - bucketName string
//   - filePath string
//   - fileType *string
func (_e *MockIClient_Expecter) NewStreamUpload(ctx interface{}, bucketName interface{}, filePath interface{}, fileType interface{}) *MockIClient_NewStreamUpload_Call {
	return &MockIClient_NewStreamUpload_Call{Call: _e.mock.On("NewStreamUpload", ctx, bucketName, filePath, fileType)}
}

func (_c *MockIClient_NewStreamUpload_Call) Run(run func(ctx context.Context, bucketName string, filePath string, fileType *string)) *MockIClient_NewStreamUpload_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(*string))
	})
	return _c
}

func (_c *MockIClient_NewStreamUpload_Call) Return(stream cloudstorage.StreamUpload, err error) *MockIClient_NewStreamUpload_Call {
	_c.Call.Return(stream, err)
	return _c
}

func (_c *MockIClient_NewStreamUpload_Call) RunAndReturn(run func(context.Context, string, string, *string) (cloudstorage.StreamUpload, error)) *MockIClient_NewStreamUpload_Call {
	_c.Call.Return(run)
	return _c
}

// Upload provides a mock function with given fields: ctx, bucketName, filePath, fileType, reader
func (_m *MockIClient) Upload(ctx context.Context, bucketName string, filePath string, fileType *string, reader io.Reader) (int64, error) {
	ret := _m.Called(ctx, bucketName, filePath, fileType, reader)

	if len(ret) == 0 {
		panic("no return value specified for Upload")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *string, io.Reader) (int64, error)); ok {
		return rf(ctx, bucketName, filePath, fileType, reader)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *string, io.Reader) int64); ok {
		r0 = rf(ctx, bucketName, filePath, fileType, reader)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, *string, io.Reader) error); ok {
		r1 = rf(ctx, bucketName, filePath, fileType, reader)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIClient_Upload_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Upload'
type MockIClient_Upload_Call struct {
	*mock.Call
}

// Upload is a helper method to define mock.On call
//   - ctx context.Context
//   - bucketName string
//   - filePath string
//   - fileType *string
//   - reader io.Reader
func (_e *MockIClient_Expecter) Upload(ctx interface{}, bucketName interface{}, filePath interface{}, fileType interface{}, reader interface{}) *MockIClient_Upload_Call {
	return &MockIClient_Upload_Call{Call: _e.mock.On("Upload", ctx, bucketName, filePath, fileType, reader)}
}

func (_c *MockIClient_Upload_Call) Run(run func(ctx context.Context, bucketName string, filePath string, fileType *string, reader io.Reader)) *MockIClient_Upload_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(*string), args[4].(io.Reader))
	})
	return _c
}

func (_c *MockIClient_Upload_Call) Return(written int64, err error) *MockIClient_Upload_Call {
	_c.Call.Return(written, err)
	return _c
}

func (_c *MockIClient_Upload_Call) RunAndReturn(run func(context.Context, string, string, *string, io.Reader) (int64, error)) *MockIClient_Upload_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIClient creates a new instance of MockIClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIClient {
	mock := &MockIClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
