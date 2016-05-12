package mocks

import "github.com/stretchr/testify/mock"

//ServicePlanInterface is a mock method for service plan
type ServicePlanInterface struct {
	mock.Mock
}

//Update is the update method mock
func (_m *ServicePlanInterface) Update(serviceGUID string) error {
	ret := _m.Called(serviceGUID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(serviceGUID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
