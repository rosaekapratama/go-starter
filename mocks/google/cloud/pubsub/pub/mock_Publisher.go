// Code generated by mockery v2.33.3. DO NOT EDIT.

package pub

import (
	context "context"

	pub "github.com/rosaekapratama/go-starter/google/cloud/pubsub/pub"
	mock "github.com/stretchr/testify/mock"
)

// MockPublisher is an autogenerated mock type for the Publisher type
type MockPublisher struct {
	mock.Mock
}

type MockPublisher_Expecter struct {
	mock *mock.Mock
}

func (_m *MockPublisher) EXPECT() *MockPublisher_Expecter {
	return &MockPublisher_Expecter{mock: &_m.Mock}
}

// BatchPublish provides a mock function with given fields: ctx, batchData, opts
func (_m *MockPublisher) BatchPublish(ctx context.Context, batchData []interface{}, opts ...pub.PublishOption) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, batchData)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []interface{}, ...pub.PublishOption) error); ok {
		r0 = rf(ctx, batchData, opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockPublisher_BatchPublish_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BatchPublish'
type MockPublisher_BatchPublish_Call struct {
	*mock.Call
}

// BatchPublish is a helper method to define mock.On call
//   - ctx context.Context
//   - batchData []interface{}
//   - opts ...pub.PublishOption
func (_e *MockPublisher_Expecter) BatchPublish(ctx interface{}, batchData interface{}, opts ...interface{}) *MockPublisher_BatchPublish_Call {
	return &MockPublisher_BatchPublish_Call{Call: _e.mock.On("BatchPublish",
		append([]interface{}{ctx, batchData}, opts...)...)}
}

func (_c *MockPublisher_BatchPublish_Call) Run(run func(ctx context.Context, batchData []interface{}, opts ...pub.PublishOption)) *MockPublisher_BatchPublish_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]pub.PublishOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(pub.PublishOption)
			}
		}
		run(args[0].(context.Context), args[1].([]interface{}), variadicArgs...)
	})
	return _c
}

func (_c *MockPublisher_BatchPublish_Call) Return(_a0 error) *MockPublisher_BatchPublish_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPublisher_BatchPublish_Call) RunAndReturn(run func(context.Context, []interface{}, ...pub.PublishOption) error) *MockPublisher_BatchPublish_Call {
	_c.Call.Return(run)
	return _c
}

// Publish provides a mock function with given fields: ctx, data, opts
func (_m *MockPublisher) Publish(ctx context.Context, data interface{}, opts ...pub.PublishOption) (string, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, data)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...pub.PublishOption) (string, error)); ok {
		return rf(ctx, data, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...pub.PublishOption) string); ok {
		r0 = rf(ctx, data, opts...)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}, ...pub.PublishOption) error); ok {
		r1 = rf(ctx, data, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockPublisher_Publish_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Publish'
type MockPublisher_Publish_Call struct {
	*mock.Call
}

// Publish is a helper method to define mock.On call
//   - ctx context.Context
//   - data interface{}
//   - opts ...pub.PublishOption
func (_e *MockPublisher_Expecter) Publish(ctx interface{}, data interface{}, opts ...interface{}) *MockPublisher_Publish_Call {
	return &MockPublisher_Publish_Call{Call: _e.mock.On("Publish",
		append([]interface{}{ctx, data}, opts...)...)}
}

func (_c *MockPublisher_Publish_Call) Run(run func(ctx context.Context, data interface{}, opts ...pub.PublishOption)) *MockPublisher_Publish_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]pub.PublishOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(pub.PublishOption)
			}
		}
		run(args[0].(context.Context), args[1].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *MockPublisher_Publish_Call) Return(serverId string, err error) *MockPublisher_Publish_Call {
	_c.Call.Return(serverId, err)
	return _c
}

func (_c *MockPublisher_Publish_Call) RunAndReturn(run func(context.Context, interface{}, ...pub.PublishOption) (string, error)) *MockPublisher_Publish_Call {
	_c.Call.Return(run)
	return _c
}

// WithAvroEncoder provides a mock function with given fields: schemaName
func (_m *MockPublisher) WithAvroEncoder(schemaName string) pub.Publisher {
	ret := _m.Called(schemaName)

	var r0 pub.Publisher
	if rf, ok := ret.Get(0).(func(string) pub.Publisher); ok {
		r0 = rf(schemaName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pub.Publisher)
		}
	}

	return r0
}

// MockPublisher_WithAvroEncoder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithAvroEncoder'
type MockPublisher_WithAvroEncoder_Call struct {
	*mock.Call
}

// WithAvroEncoder is a helper method to define mock.On call
//   - schemaName string
func (_e *MockPublisher_Expecter) WithAvroEncoder(schemaName interface{}) *MockPublisher_WithAvroEncoder_Call {
	return &MockPublisher_WithAvroEncoder_Call{Call: _e.mock.On("WithAvroEncoder", schemaName)}
}

func (_c *MockPublisher_WithAvroEncoder_Call) Run(run func(schemaName string)) *MockPublisher_WithAvroEncoder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockPublisher_WithAvroEncoder_Call) Return(_a0 pub.Publisher) *MockPublisher_WithAvroEncoder_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPublisher_WithAvroEncoder_Call) RunAndReturn(run func(string) pub.Publisher) *MockPublisher_WithAvroEncoder_Call {
	_c.Call.Return(run)
	return _c
}

// WithJsonEncoder provides a mock function with given fields:
func (_m *MockPublisher) WithJsonEncoder() pub.Publisher {
	ret := _m.Called()

	var r0 pub.Publisher
	if rf, ok := ret.Get(0).(func() pub.Publisher); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pub.Publisher)
		}
	}

	return r0
}

// MockPublisher_WithJsonEncoder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithJsonEncoder'
type MockPublisher_WithJsonEncoder_Call struct {
	*mock.Call
}

// WithJsonEncoder is a helper method to define mock.On call
func (_e *MockPublisher_Expecter) WithJsonEncoder() *MockPublisher_WithJsonEncoder_Call {
	return &MockPublisher_WithJsonEncoder_Call{Call: _e.mock.On("WithJsonEncoder")}
}

func (_c *MockPublisher_WithJsonEncoder_Call) Run(run func()) *MockPublisher_WithJsonEncoder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockPublisher_WithJsonEncoder_Call) Return(_a0 pub.Publisher) *MockPublisher_WithJsonEncoder_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPublisher_WithJsonEncoder_Call) RunAndReturn(run func() pub.Publisher) *MockPublisher_WithJsonEncoder_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockPublisher creates a new instance of MockPublisher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPublisher(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPublisher {
	mock := &MockPublisher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
