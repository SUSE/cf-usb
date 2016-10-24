package mocks

import "github.com/stretchr/testify/mock"

type CSM struct {
	mock.Mock
}

// Login provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *CSM) Login(_a0 string, _a1 string, _a2 string, _a3 bool) error {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, bool) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateWorkspace provides a mock function with given fields: _a0
func (_m *CSM) CreateWorkspace(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WorkspaceExists provides a mock function with given fields: _a0
func (_m *CSM) WorkspaceExists(_a0 string) (bool, bool, error) {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(_a0)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DeleteWorkspace provides a mock function with given fields: _a0
func (_m *CSM) DeleteWorkspace(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateConnection provides a mock function with given fields: _a0, _a1
func (_m *CSM) CreateConnection(_a0 string, _a1 string) (interface{}, error) {
	ret := _m.Called(_a0, _a1)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string, string) interface{}); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ConnectionExists provides a mock function with given fields: _a0, _a1
func (_m *CSM) ConnectionExists(_a0 string, _a1 string) (bool, bool, error) {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string, string) bool); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string) error); ok {
		r2 = rf(_a0, _a1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DeleteConnection provides a mock function with given fields: _a0, _a1
func (_m *CSM) DeleteConnection(_a0 string, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetStatus provides a mock function with given fields:
func (_m *CSM) GetStatus() (string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
