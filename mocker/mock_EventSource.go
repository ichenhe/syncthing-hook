// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocker

import (
	domain "github.com/ichenhe/syncthing-hook/domain"
	mock "github.com/stretchr/testify/mock"
)

// MockEventSource is an autogenerated mock type for the EventSource type
type MockEventSource struct {
	mock.Mock
}

type MockEventSource_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEventSource) EXPECT() *MockEventSource_Expecter {
	return &MockEventSource_Expecter{mock: &_m.Mock}
}

// Subscribe provides a mock function with given fields: eventType, params, hookDef
func (_m *MockEventSource) Subscribe(eventType domain.EventType, params *domain.HookParameters, hookDef *domain.HookDefinition) (<-chan *domain.Event, error) {
	ret := _m.Called(eventType, params, hookDef)

	var r0 <-chan *domain.Event
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.EventType, *domain.HookParameters, *domain.HookDefinition) (<-chan *domain.Event, error)); ok {
		return rf(eventType, params, hookDef)
	}
	if rf, ok := ret.Get(0).(func(domain.EventType, *domain.HookParameters, *domain.HookDefinition) <-chan *domain.Event); ok {
		r0 = rf(eventType, params, hookDef)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan *domain.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(domain.EventType, *domain.HookParameters, *domain.HookDefinition) error); ok {
		r1 = rf(eventType, params, hookDef)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockEventSource_Subscribe_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Subscribe'
type MockEventSource_Subscribe_Call struct {
	*mock.Call
}

// Subscribe is a helper method to define mock.On call
//   - eventType domain.EventType
//   - params *domain.HookParameters
//   - hookDef *domain.HookDefinition
func (_e *MockEventSource_Expecter) Subscribe(eventType interface{}, params interface{}, hookDef interface{}) *MockEventSource_Subscribe_Call {
	return &MockEventSource_Subscribe_Call{Call: _e.mock.On("Subscribe", eventType, params, hookDef)}
}

func (_c *MockEventSource_Subscribe_Call) Run(run func(eventType domain.EventType, params *domain.HookParameters, hookDef *domain.HookDefinition)) *MockEventSource_Subscribe_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(domain.EventType), args[1].(*domain.HookParameters), args[2].(*domain.HookDefinition))
	})
	return _c
}

func (_c *MockEventSource_Subscribe_Call) Return(_a0 <-chan *domain.Event, _a1 error) *MockEventSource_Subscribe_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockEventSource_Subscribe_Call) RunAndReturn(run func(domain.EventType, *domain.HookParameters, *domain.HookDefinition) (<-chan *domain.Event, error)) *MockEventSource_Subscribe_Call {
	_c.Call.Return(run)
	return _c
}

// Unsubscribe provides a mock function with given fields: eventCh
func (_m *MockEventSource) Unsubscribe(eventCh <-chan *domain.Event) {
	_m.Called(eventCh)
}

// MockEventSource_Unsubscribe_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Unsubscribe'
type MockEventSource_Unsubscribe_Call struct {
	*mock.Call
}

// Unsubscribe is a helper method to define mock.On call
//   - eventCh <-chan *domain.Event
func (_e *MockEventSource_Expecter) Unsubscribe(eventCh interface{}) *MockEventSource_Unsubscribe_Call {
	return &MockEventSource_Unsubscribe_Call{Call: _e.mock.On("Unsubscribe", eventCh)}
}

func (_c *MockEventSource_Unsubscribe_Call) Run(run func(eventCh <-chan *domain.Event)) *MockEventSource_Unsubscribe_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(<-chan *domain.Event))
	})
	return _c
}

func (_c *MockEventSource_Unsubscribe_Call) Return() *MockEventSource_Unsubscribe_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockEventSource_Unsubscribe_Call) RunAndReturn(run func(<-chan *domain.Event)) *MockEventSource_Unsubscribe_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockEventSource creates a new instance of MockEventSource. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEventSource(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEventSource {
	mock := &MockEventSource{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}