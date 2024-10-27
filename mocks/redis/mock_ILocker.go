// Code generated by mockery v2.46.3. DO NOT EDIT.

package redis

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	redislock "github.com/bsm/redislock"

	time "time"
)

// MockILocker is an autogenerated mock type for the ILocker type
type MockILocker struct {
	mock.Mock
}

type MockILocker_Expecter struct {
	mock *mock.Mock
}

func (_m *MockILocker) EXPECT() *MockILocker_Expecter {
	return &MockILocker_Expecter{mock: &_m.Mock}
}

// Obtain provides a mock function with given fields: ctx, key, ttl, opt
func (_m *MockILocker) Obtain(ctx context.Context, key string, ttl time.Duration, opt *redislock.Options) (*redislock.Lock, error) {
	ret := _m.Called(ctx, key, ttl, opt)

	if len(ret) == 0 {
		panic("no return value specified for Obtain")
	}

	var r0 *redislock.Lock
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Duration, *redislock.Options) (*redislock.Lock, error)); ok {
		return rf(ctx, key, ttl, opt)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Duration, *redislock.Options) *redislock.Lock); ok {
		r0 = rf(ctx, key, ttl, opt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*redislock.Lock)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, time.Duration, *redislock.Options) error); ok {
		r1 = rf(ctx, key, ttl, opt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockILocker_Obtain_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Obtain'
type MockILocker_Obtain_Call struct {
	*mock.Call
}

// Obtain is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - ttl time.Duration
//   - opt *redislock.Options
func (_e *MockILocker_Expecter) Obtain(ctx interface{}, key interface{}, ttl interface{}, opt interface{}) *MockILocker_Obtain_Call {
	return &MockILocker_Obtain_Call{Call: _e.mock.On("Obtain", ctx, key, ttl, opt)}
}

func (_c *MockILocker_Obtain_Call) Run(run func(ctx context.Context, key string, ttl time.Duration, opt *redislock.Options)) *MockILocker_Obtain_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(time.Duration), args[3].(*redislock.Options))
	})
	return _c
}

func (_c *MockILocker_Obtain_Call) Return(_a0 *redislock.Lock, _a1 error) *MockILocker_Obtain_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockILocker_Obtain_Call) RunAndReturn(run func(context.Context, string, time.Duration, *redislock.Options) (*redislock.Lock, error)) *MockILocker_Obtain_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockILocker creates a new instance of MockILocker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockILocker(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockILocker {
	mock := &MockILocker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
