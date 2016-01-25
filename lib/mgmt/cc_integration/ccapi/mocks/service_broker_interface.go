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
func (_m *ServiceBrokerInterface) Update(serviceBrokerGuid string, name string, url string, username string, password string) error {
	ret := _m.Called(serviceBrokerGuid, name, url, username, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, string, string) error); ok {
		r0 = rf(serviceBrokerGuid, name, url, username, password)
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
func (_m *ServiceBrokerInterface) GetServiceBrokerGuidByName(name string) (string, error) {
	ret := _m.Called(name)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceBrokerInterface) CheckServiceNameExists(name string) bool {
	ret := _m.Called(name)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
