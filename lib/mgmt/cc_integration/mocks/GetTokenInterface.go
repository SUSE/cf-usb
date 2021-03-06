package mocks

import "github.com/stretchr/testify/mock"

//GetTokenInterface is a mock for token interface
type GetTokenInterface struct {
	mock.Mock
}

//GetToken mocks GetToken function
func (_m *GetTokenInterface) GetToken() (string, error) {
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
