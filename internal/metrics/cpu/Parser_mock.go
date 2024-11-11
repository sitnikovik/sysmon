// Code generated by mockery. DO NOT EDIT.

package cpu

import mock "github.com/stretchr/testify/mock"

// MockParser is an autogenerated mock type for the Parser type
type MockParser struct {
	mock.Mock
}

type MockParser_Expecter struct {
	mock *mock.Mock
}

func (_m *MockParser) EXPECT() *MockParser_Expecter {
	return &MockParser_Expecter{mock: &_m.Mock}
}

// Parse provides a mock function with given fields:
func (_m *MockParser) Parse() (CpuStats, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Parse")
	}

	var r0 CpuStats
	var r1 error
	if rf, ok := ret.Get(0).(func() (CpuStats, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() CpuStats); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(CpuStats)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockParser_Parse_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Parse'
type MockParser_Parse_Call struct {
	*mock.Call
}

// Parse is a helper method to define mock.On call
func (_e *MockParser_Expecter) Parse() *MockParser_Parse_Call {
	return &MockParser_Parse_Call{Call: _e.mock.On("Parse")}
}

func (_c *MockParser_Parse_Call) Run(run func()) *MockParser_Parse_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockParser_Parse_Call) Return(_a0 CpuStats, _a1 error) *MockParser_Parse_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockParser_Parse_Call) RunAndReturn(run func() (CpuStats, error)) *MockParser_Parse_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockParser creates a new instance of MockParser. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockParser(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockParser {
	mock := &MockParser{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
