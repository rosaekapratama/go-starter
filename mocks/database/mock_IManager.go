// Code generated by mockery v2.46.3. DO NOT EDIT.

package database

import (
	context "context"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	sql "database/sql"
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

// Begin provides a mock function with given fields: ctx, connectionId
func (_m *MockIManager) Begin(ctx context.Context, connectionId string) (*gorm.DB, error) {
	ret := _m.Called(ctx, connectionId)

	if len(ret) == 0 {
		panic("no return value specified for Begin")
	}

	var r0 *gorm.DB
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*gorm.DB, error)); ok {
		return rf(ctx, connectionId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *gorm.DB); ok {
		r0 = rf(ctx, connectionId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, connectionId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIManager_Begin_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Begin'
type MockIManager_Begin_Call struct {
	*mock.Call
}

// Begin is a helper method to define mock.On call
//   - ctx context.Context
//   - connectionId string
func (_e *MockIManager_Expecter) Begin(ctx interface{}, connectionId interface{}) *MockIManager_Begin_Call {
	return &MockIManager_Begin_Call{Call: _e.mock.On("Begin", ctx, connectionId)}
}

func (_c *MockIManager_Begin_Call) Run(run func(ctx context.Context, connectionId string)) *MockIManager_Begin_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockIManager_Begin_Call) Return(_a0 *gorm.DB, _a1 error) *MockIManager_Begin_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIManager_Begin_Call) RunAndReturn(run func(context.Context, string) (*gorm.DB, error)) *MockIManager_Begin_Call {
	_c.Call.Return(run)
	return _c
}

// DB provides a mock function with given fields: ctx, connectionId
func (_m *MockIManager) DB(ctx context.Context, connectionId string) (*gorm.DB, *sql.DB, error) {
	ret := _m.Called(ctx, connectionId)

	if len(ret) == 0 {
		panic("no return value specified for DB")
	}

	var r0 *gorm.DB
	var r1 *sql.DB
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*gorm.DB, *sql.DB, error)); ok {
		return rf(ctx, connectionId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *gorm.DB); ok {
		r0 = rf(ctx, connectionId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) *sql.DB); ok {
		r1 = rf(ctx, connectionId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*sql.DB)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, connectionId)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockIManager_DB_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DB'
type MockIManager_DB_Call struct {
	*mock.Call
}

// DB is a helper method to define mock.On call
//   - ctx context.Context
//   - connectionId string
func (_e *MockIManager_Expecter) DB(ctx interface{}, connectionId interface{}) *MockIManager_DB_Call {
	return &MockIManager_DB_Call{Call: _e.mock.On("DB", ctx, connectionId)}
}

func (_c *MockIManager_DB_Call) Run(run func(ctx context.Context, connectionId string)) *MockIManager_DB_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockIManager_DB_Call) Return(_a0 *gorm.DB, _a1 *sql.DB, _a2 error) *MockIManager_DB_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockIManager_DB_Call) RunAndReturn(run func(context.Context, string) (*gorm.DB, *sql.DB, error)) *MockIManager_DB_Call {
	_c.Call.Return(run)
	return _c
}

// GetConnectionIds provides a mock function with given fields:
func (_m *MockIManager) GetConnectionIds() []string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetConnectionIds")
	}

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// MockIManager_GetConnectionIds_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConnectionIds'
type MockIManager_GetConnectionIds_Call struct {
	*mock.Call
}

// GetConnectionIds is a helper method to define mock.On call
func (_e *MockIManager_Expecter) GetConnectionIds() *MockIManager_GetConnectionIds_Call {
	return &MockIManager_GetConnectionIds_Call{Call: _e.mock.On("GetConnectionIds")}
}

func (_c *MockIManager_GetConnectionIds_Call) Run(run func()) *MockIManager_GetConnectionIds_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockIManager_GetConnectionIds_Call) Return(_a0 []string) *MockIManager_GetConnectionIds_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIManager_GetConnectionIds_Call) RunAndReturn(run func() []string) *MockIManager_GetConnectionIds_Call {
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
