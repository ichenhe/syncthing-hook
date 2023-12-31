// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocker

import (
	domain "github.com/ichenhe/syncthing-hook/domain"
	mock "github.com/stretchr/testify/mock"
)

// MockEventHandler is an autogenerated mock type for the EventHandler type
type MockEventHandler struct {
	mock.Mock
}

type MockEventHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEventHandler) EXPECT() *MockEventHandler_Expecter {
	return &MockEventHandler_Expecter{mock: &_m.Mock}
}

// Destroy provides a mock function with given fields:
func (_m *MockEventHandler) Destroy() {
	_m.Called()
}

// MockEventHandler_Destroy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Destroy'
type MockEventHandler_Destroy_Call struct {
	*mock.Call
}

// Destroy is a helper method to define mock.On call
func (_e *MockEventHandler_Expecter) Destroy() *MockEventHandler_Destroy_Call {
	return &MockEventHandler_Destroy_Call{Call: _e.mock.On("Destroy")}
}

func (_c *MockEventHandler_Destroy_Call) Run(run func()) *MockEventHandler_Destroy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEventHandler_Destroy_Call) Return() *MockEventHandler_Destroy_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockEventHandler_Destroy_Call) RunAndReturn(run func()) *MockEventHandler_Destroy_Call {
	_c.Call.Return(run)
	return _c
}

// GetNext provides a mock function with given fields:
func (_m *MockEventHandler) GetNext() domain.EventHandler {
	ret := _m.Called()

	var r0 domain.EventHandler
	if rf, ok := ret.Get(0).(func() domain.EventHandler); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.EventHandler)
		}
	}

	return r0
}

// MockEventHandler_GetNext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetNext'
type MockEventHandler_GetNext_Call struct {
	*mock.Call
}

// GetNext is a helper method to define mock.On call
func (_e *MockEventHandler_Expecter) GetNext() *MockEventHandler_GetNext_Call {
	return &MockEventHandler_GetNext_Call{Call: _e.mock.On("GetNext")}
}

func (_c *MockEventHandler_GetNext_Call) Run(run func()) *MockEventHandler_GetNext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEventHandler_GetNext_Call) Return(_a0 domain.EventHandler) *MockEventHandler_GetNext_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockEventHandler_GetNext_Call) RunAndReturn(run func() domain.EventHandler) *MockEventHandler_GetNext_Call {
	_c.Call.Return(run)
	return _c
}

// Handle provides a mock function with given fields: event
func (_m *MockEventHandler) Handle(event *domain.Event) {
	_m.Called(event)
}

// MockEventHandler_Handle_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Handle'
type MockEventHandler_Handle_Call struct {
	*mock.Call
}

// Handle is a helper method to define mock.On call
//   - event *domain.Event
func (_e *MockEventHandler_Expecter) Handle(event interface{}) *MockEventHandler_Handle_Call {
	return &MockEventHandler_Handle_Call{Call: _e.mock.On("Handle", event)}
}

func (_c *MockEventHandler_Handle_Call) Run(run func(event *domain.Event)) *MockEventHandler_Handle_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*domain.Event))
	})
	return _c
}

func (_c *MockEventHandler_Handle_Call) Return() *MockEventHandler_Handle_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockEventHandler_Handle_Call) RunAndReturn(run func(*domain.Event)) *MockEventHandler_Handle_Call {
	_c.Call.Return(run)
	return _c
}

// SetNext provides a mock function with given fields: next
func (_m *MockEventHandler) SetNext(next domain.EventHandler) {
	_m.Called(next)
}

// MockEventHandler_SetNext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetNext'
type MockEventHandler_SetNext_Call struct {
	*mock.Call
}

// SetNext is a helper method to define mock.On call
//   - next domain.EventHandler
func (_e *MockEventHandler_Expecter) SetNext(next interface{}) *MockEventHandler_SetNext_Call {
	return &MockEventHandler_SetNext_Call{Call: _e.mock.On("SetNext", next)}
}

func (_c *MockEventHandler_SetNext_Call) Run(run func(next domain.EventHandler)) *MockEventHandler_SetNext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(domain.EventHandler))
	})
	return _c
}

func (_c *MockEventHandler_SetNext_Call) Return() *MockEventHandler_SetNext_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockEventHandler_SetNext_Call) RunAndReturn(run func(domain.EventHandler)) *MockEventHandler_SetNext_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockEventHandler creates a new instance of MockEventHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEventHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEventHandler {
	mock := &MockEventHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
