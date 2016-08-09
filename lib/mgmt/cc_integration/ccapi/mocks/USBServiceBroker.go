package mocks

import "github.com/stretchr/testify/mock"

type USBServiceBroker struct {
	mock.Mock
}

// Create provides a mock function with given fields: name, url, username, password
func (_m *USBServiceBroker) Create(name string, url string, username string, password string) error {
	ret := _m.Called(name, url, username, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, string) error); ok {
		r0 = rf(name, url, username, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: name
func (_m *USBServiceBroker) Delete(name string) error {
	ret := _m.Called(name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: serviceBrokerGUID, name, url, username, password
func (_m *USBServiceBroker) Update(serviceBrokerGUID string, name string, url string, username string, password string) error {
	ret := _m.Called(serviceBrokerGUID, name, url, username, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, string, string) error); ok {
		r0 = rf(serviceBrokerGUID, name, url, username, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EnableServiceAccess provides a mock function with given fields: serviceID
func (_m *USBServiceBroker) EnableServiceAccess(serviceID string) error {
	ret := _m.Called(serviceID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(serviceID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetServiceBrokerGUIDByName provides a mock function with given fields: name
func (_m *USBServiceBroker) GetServiceBrokerGUIDByName(name string) (string, error) {
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

// CheckServiceNameExists provides a mock function with given fields: name
func (_m *USBServiceBroker) CheckServiceNameExists(name string) (bool, error) {
	ret := _m.Called(name)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckServiceInstancesExist provides a mock function with given fields: serviceName
func (_m *USBServiceBroker) CheckServiceInstancesExist(serviceName string) bool {
	ret := _m.Called(serviceName)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(serviceName)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
