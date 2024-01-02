// Code generated by mockery v2.33.3. DO NOT EDIT.

package grpcclient

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// MockIManager is an autogenerated mock type for the IManager type
type MockIManager struct {
	mock.Mock
}

type MockIManager_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIManager) EXPECT() *MockIManager_Expecter {
	return &MockIManager_Expecter{mock: &_m.Mock}
}

// GetConn provides a mock function with given fields: ctx, connId
func (_m *MockIManager) GetConn(ctx context.Context, connId string) (*grpc.ClientConn, error) {
	ret := _m.Called(ctx, connId)

	var r0 *grpc.ClientConn
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*grpc.ClientConn, error)); ok {
		return rf(ctx, connId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *grpc.ClientConn); ok {
		r0 = rf(ctx, connId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpc.ClientConn)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, connId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIManager_GetConn_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConn'
type MockIManager_GetConn_Call struct {
	*mock.Call
}

// GetConn is a helper method to define mock.On call
//   - ctx context.Context
//   - connId string
func (_e *MockIManager_Expecter) GetConn(ctx interface{}, connId interface{}) *MockIManager_GetConn_Call {
	return &MockIManager_GetConn_Call{Call: _e.mock.On("GetConn", ctx, connId)}
}

func (_c *MockIManager_GetConn_Call) Run(run func(ctx context.Context, connId string)) *MockIManager_GetConn_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockIManager_GetConn_Call) Return(conn *grpc.ClientConn, err error) *MockIManager_GetConn_Call {
	_c.Call.Return(conn, err)
	return _c
}

func (_c *MockIManager_GetConn_Call) RunAndReturn(run func(context.Context, string) (*grpc.ClientConn, error)) *MockIManager_GetConn_Call {
	_c.Call.Return(run)
	return _c
}

// InitConn provides a mock function with given fields: ctx, connId, address, opts
func (_m *MockIManager) InitConn(ctx context.Context, connId string, address string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, connId, address)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *grpc.ClientConn
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, ...grpc.DialOption) (*grpc.ClientConn, error)); ok {
		return rf(ctx, connId, address, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, ...grpc.DialOption) *grpc.ClientConn); ok {
		r0 = rf(ctx, connId, address, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpc.ClientConn)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, ...grpc.DialOption) error); ok {
		r1 = rf(ctx, connId, address, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIManager_InitConn_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InitConn'
type MockIManager_InitConn_Call struct {
	*mock.Call
}

// InitConn is a helper method to define mock.On call
//   - ctx context.Context
//   - connId string
//   - address string
//   - opts ...grpc.DialOption
func (_e *MockIManager_Expecter) InitConn(ctx interface{}, connId interface{}, address interface{}, opts ...interface{}) *MockIManager_InitConn_Call {
	return &MockIManager_InitConn_Call{Call: _e.mock.On("InitConn",
		append([]interface{}{ctx, connId, address}, opts...)...)}
}

func (_c *MockIManager_InitConn_Call) Run(run func(ctx context.Context, connId string, address string, opts ...grpc.DialOption)) *MockIManager_InitConn_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.DialOption, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.DialOption)
			}
		}
		run(args[0].(context.Context), args[1].(string), args[2].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockIManager_InitConn_Call) Return(conn *grpc.ClientConn, err error) *MockIManager_InitConn_Call {
	_c.Call.Return(conn, err)
	return _c
}

func (_c *MockIManager_InitConn_Call) RunAndReturn(run func(context.Context, string, string, ...grpc.DialOption) (*grpc.ClientConn, error)) *MockIManager_InitConn_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIManager creates a new instance of MockIManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIManager {
	mock := &MockIManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
