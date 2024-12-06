// Code generated by mockery v2.49.2. DO NOT EDIT.

//go:build linux && test

package nvidia

import mock "github.com/stretchr/testify/mock"

// mockCollector is an autogenerated mock type for the Collector type
type mockCollector struct {
	mock.Mock
}

type mockCollector_Expecter struct {
	mock *mock.Mock
}

func (_m *mockCollector) EXPECT() *mockCollector_Expecter {
	return &mockCollector_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with no fields
func (_m *mockCollector) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockCollector_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type mockCollector_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *mockCollector_Expecter) Close() *mockCollector_Close_Call {
	return &mockCollector_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *mockCollector_Close_Call) Run(run func()) *mockCollector_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockCollector_Close_Call) Return(_a0 error) *mockCollector_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCollector_Close_Call) RunAndReturn(run func() error) *mockCollector_Close_Call {
	_c.Call.Return(run)
	return _c
}

// Collect provides a mock function with no fields
func (_m *mockCollector) Collect() ([]Metric, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Collect")
	}

	var r0 []Metric
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]Metric, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []Metric); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Metric)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockCollector_Collect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Collect'
type mockCollector_Collect_Call struct {
	*mock.Call
}

// Collect is a helper method to define mock.On call
func (_e *mockCollector_Expecter) Collect() *mockCollector_Collect_Call {
	return &mockCollector_Collect_Call{Call: _e.mock.On("Collect")}
}

func (_c *mockCollector_Collect_Call) Run(run func()) *mockCollector_Collect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockCollector_Collect_Call) Return(_a0 []Metric, _a1 error) *mockCollector_Collect_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockCollector_Collect_Call) RunAndReturn(run func() ([]Metric, error)) *mockCollector_Collect_Call {
	_c.Call.Return(run)
	return _c
}

// Name provides a mock function with no fields
func (_m *mockCollector) Name() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Name")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// mockCollector_Name_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Name'
type mockCollector_Name_Call struct {
	*mock.Call
}

// Name is a helper method to define mock.On call
func (_e *mockCollector_Expecter) Name() *mockCollector_Name_Call {
	return &mockCollector_Name_Call{Call: _e.mock.On("Name")}
}

func (_c *mockCollector_Name_Call) Run(run func()) *mockCollector_Name_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockCollector_Name_Call) Return(_a0 string) *mockCollector_Name_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCollector_Name_Call) RunAndReturn(run func() string) *mockCollector_Name_Call {
	_c.Call.Return(run)
	return _c
}

// newMockCollector creates a new instance of mockCollector. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockCollector(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockCollector {
	mock := &mockCollector{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
