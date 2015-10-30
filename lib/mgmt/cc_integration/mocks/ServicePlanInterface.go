package mocks

import "github.com/stretchr/testify/mock"

type ServicePlanInterface struct {
	mock.Mock
}

func (_m *ServicePlanInterface) Update(serviceGuid string) error {
	ret := _m.Called(serviceGuid)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(serviceGuid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
