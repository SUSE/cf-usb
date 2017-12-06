package mocks

import "github.com/SUSE/cf-usb/lib/config"
import "github.com/stretchr/testify/mock"

import "github.com/SUSE/cf-usb/lib/brokermodel"

type Provider struct {
	mock.Mock
}

// InitializeConfiguration does nothing for the mock provider
func (_m *Provider) InitializeConfiguration() error {
	return nil
}

// LoadConfiguration provides a mock function with given fields:
func (_m *Provider) LoadConfiguration() (*config.Config, error) {
	ret := _m.Called()

	var r0 *config.Config
	if rf, ok := ret.Get(0).(func() *config.Config); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.Config)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveConfiguration provides a mock function with given fields: config, overwrite
func (_m *Provider) SaveConfiguration(config config.Config, overwrite bool) error {
	return nil
}

// LoadDriverInstance provides a mock function with given fields: driverInstanceID
func (_m *Provider) LoadDriverInstance(driverInstanceID string) (*config.Instance, error) {
	ret := _m.Called(driverInstanceID)

	var r0 *config.Instance
	if rf, ok := ret.Get(0).(func(string) *config.Instance); ok {
		r0 = rf(driverInstanceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.Instance)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(driverInstanceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUaaAuthConfig provides a mock function with given fields:
func (_m *Provider) GetUaaAuthConfig() (*config.UaaAuth, error) {
	ret := _m.Called()

	var r0 *config.UaaAuth
	if rf, ok := ret.Get(0).(func() *config.UaaAuth); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.UaaAuth)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetInstance provides a mock function with given fields: instanceid, driverInstance
func (_m *Provider) SetInstance(instanceid string, driverInstance config.Instance) error {
	ret := _m.Called(instanceid, driverInstance)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, config.Instance) error); ok {
		r0 = rf(instanceid, driverInstance)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetInstance provides a mock function with given fields: instanceid
func (_m *Provider) GetInstance(instanceid string) (*config.Instance, string, error) {
	ret := _m.Called(instanceid)

	var r0 *config.Instance
	if rf, ok := ret.Get(0).(func(string) *config.Instance); ok {
		r0 = rf(instanceid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.Instance)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(string) string); ok {
		r1 = rf(instanceid)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(instanceid)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DeleteInstance provides a mock function with given fields: instanceid
func (_m *Provider) DeleteInstance(instanceid string) error {
	ret := _m.Called(instanceid)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(instanceid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetService provides a mock function with given fields: instanceid, service
func (_m *Provider) SetService(instanceid string, service brokermodel.CatalogService) error {
	ret := _m.Called(instanceid, service)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, brokermodel.CatalogService) error); ok {
		r0 = rf(instanceid, service)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetService provides a mock function with given fields: serviceid
func (_m *Provider) GetService(serviceid string) (*brokermodel.CatalogService, string, error) {
	ret := _m.Called(serviceid)

	var r0 *brokermodel.CatalogService
	if rf, ok := ret.Get(0).(func(string) *brokermodel.CatalogService); ok {
		r0 = rf(serviceid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*brokermodel.CatalogService)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(string) string); ok {
		r1 = rf(serviceid)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(serviceid)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DeleteService provides a mock function with given fields: instanceid
func (_m *Provider) DeleteService(instanceid string) error {
	ret := _m.Called(instanceid)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(instanceid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetDial provides a mock function with given fields: instanceid, dialid, dial
func (_m *Provider) SetDial(instanceid string, dialid string, dial config.Dial) error {
	ret := _m.Called(instanceid, dialid, dial)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, config.Dial) error); ok {
		r0 = rf(instanceid, dialid, dial)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetDial provides a mock function with given fields: dialid
func (_m *Provider) GetDial(dialid string) (*config.Dial, string, error) {
	ret := _m.Called(dialid)

	var r0 *config.Dial
	if rf, ok := ret.Get(0).(func(string) *config.Dial); ok {
		r0 = rf(dialid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.Dial)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(string) string); ok {
		r1 = rf(dialid)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(dialid)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DeleteDial provides a mock function with given fields: dialid
func (_m *Provider) DeleteDial(dialid string) error {
	ret := _m.Called(dialid)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(dialid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InstanceNameExists provides a mock function with given fields: driverInstanceName
func (_m *Provider) InstanceNameExists(driverInstanceName string) (bool, error) {
	ret := _m.Called(driverInstanceName)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(driverInstanceName)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(driverInstanceName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPlan provides a mock function with given fields: plandid
func (_m *Provider) GetPlan(plandid string) (*brokermodel.Plan, string, string, error) {
	ret := _m.Called(plandid)

	var r0 *brokermodel.Plan
	if rf, ok := ret.Get(0).(func(string) *brokermodel.Plan); ok {
		r0 = rf(plandid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*brokermodel.Plan)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(string) string); ok {
		r1 = rf(plandid)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 string
	if rf, ok := ret.Get(2).(func(string) string); ok {
		r2 = rf(plandid)
	} else {
		r2 = ret.Get(2).(string)
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(string) error); ok {
		r3 = rf(plandid)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}
