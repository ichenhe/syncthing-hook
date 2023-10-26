// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocker

import mock "github.com/stretchr/testify/mock"

// MockTimeProvider is an autogenerated mock type for the TimeProvider type
type MockTimeProvider struct {
	mock.Mock
}

type MockTimeProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTimeProvider) EXPECT() *MockTimeProvider_Expecter {
	return &MockTimeProvider_Expecter{mock: &_m.Mock}
}

// NowUnixMilli provides a mock function with given fields:
func (_m *MockTimeProvider) NowUnixMilli() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// MockTimeProvider_NowUnixMilli_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NowUnixMilli'
type MockTimeProvider_NowUnixMilli_Call struct {
	*mock.Call
}

// NowUnixMilli is a helper method to define mock.On call
func (_e *MockTimeProvider_Expecter) NowUnixMilli() *MockTimeProvider_NowUnixMilli_Call {
	return &MockTimeProvider_NowUnixMilli_Call{Call: _e.mock.On("NowUnixMilli")}
}

func (_c *MockTimeProvider_NowUnixMilli_Call) Run(run func()) *MockTimeProvider_NowUnixMilli_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockTimeProvider_NowUnixMilli_Call) Return(_a0 int64) *MockTimeProvider_NowUnixMilli_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTimeProvider_NowUnixMilli_Call) RunAndReturn(run func() int64) *MockTimeProvider_NowUnixMilli_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTimeProvider creates a new instance of MockTimeProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTimeProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTimeProvider {
	mock := &MockTimeProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
