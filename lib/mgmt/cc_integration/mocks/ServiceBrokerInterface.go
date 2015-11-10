package mocks

import "github.com/stretchr/testify/mock"

type ServiceBrokerInterface struct {
	mock.Mock
}

func (_m *ServiceBrokerInterface) Create(name string, url string, username string, password string) error {
	ret := _m.Called(name, url, username, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, string) error); ok {
		r0 = rf(name, url, username, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceBrokerInterface) Update(name string, url string, username string, password string) error {
	ret := _m.Called(name, url, username, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, string) error); ok {
		r0 = rf(name, url, username, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceBrokerInterface) EnableServiceAccess(serviceId string) error {
	ret := _m.Called(serviceId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(serviceId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
