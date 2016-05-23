package mocks

import "github.com/stretchr/testify/mock"

//GetInfoInterface is a mock for Info Interface
type GetInfoInterface struct {
	mock.Mock
}

//GetTokenEndpoint mocks geting the endpoint of the token
func (_m *GetInfoInterface) GetTokenEndpoint() (string, error) {
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
