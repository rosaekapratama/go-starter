// Code generated by mockery v2.33.3. DO NOT EDIT.

package restclient

import (
	context "context"

	restclient "github.com/rosaekapratama/go-starter/transport/restclient"
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

// GetDefaultClient provides a mock function with given fields:
func (_m *MockIManager) GetDefaultClient() *restclient.Client {
	ret := _m.Called()

	var r0 *restclient.Client
	if rf, ok := ret.Get(0).(func() *restclient.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*restclient.Client)
		}
	}

	return r0
}

// MockIManager_GetDefaultClient_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDefaultClient'
type MockIManager_GetDefaultClient_Call struct {
	*mock.Call
}

// GetDefaultClient is a helper method to define mock.On call
func (_e *MockIManager_Expecter) GetDefaultClient() *MockIManager_GetDefaultClient_Call {
	return &MockIManager_GetDefaultClient_Call{Call: _e.mock.On("GetDefaultClient")}
}

func (_c *MockIManager_GetDefaultClient_Call) Run(run func()) *MockIManager_GetDefaultClient_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockIManager_GetDefaultClient_Call) Return(_a0 *restclient.Client) *MockIManager_GetDefaultClient_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIManager_GetDefaultClient_Call) RunAndReturn(run func() *restclient.Client) *MockIManager_GetDefaultClient_Call {
	_c.Call.Return(run)
	return _c
}

// NewClient provides a mock function with given fields: ctx, opts
func (_m *MockIManager) NewClient(ctx context.Context, opts ...restclient.ClientOption) (*restclient.Client, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *restclient.Client
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...restclient.ClientOption) (*restclient.Client, error)); ok {
		return rf(ctx, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...restclient.ClientOption) *restclient.Client); ok {
		r0 = rf(ctx, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*restclient.Client)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...restclient.ClientOption) error); ok {
		r1 = rf(ctx, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIManager_NewClient_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NewClient'
type MockIManager_NewClient_Call struct {
	*mock.Call
}

// NewClient is a helper method to define mock.On call
//   - ctx context.Context
//   - opts ...restclient.ClientOption
func (_e *MockIManager_Expecter) NewClient(ctx interface{}, opts ...interface{}) *MockIManager_NewClient_Call {
	return &MockIManager_NewClient_Call{Call: _e.mock.On("NewClient",
		append([]interface{}{ctx}, opts...)...)}
}

func (_c *MockIManager_NewClient_Call) Run(run func(ctx context.Context, opts ...restclient.ClientOption)) *MockIManager_NewClient_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]restclient.ClientOption, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(restclient.ClientOption)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *MockIManager_NewClient_Call) Return(_a0 *restclient.Client, _a1 error) *MockIManager_NewClient_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIManager_NewClient_Call) RunAndReturn(run func(context.Context, ...restclient.ClientOption) (*restclient.Client, error)) *MockIManager_NewClient_Call {
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
