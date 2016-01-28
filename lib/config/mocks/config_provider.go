package mocks

import "github.com/hpcloud/cf-usb/lib/config"
import "github.com/stretchr/testify/mock"

import "github.com/frodenas/brokerapi"

type ConfigProvider struct {
	mock.Mock
}

func (_m *ConfigProvider) LoadConfiguration() (*config.Config, error) {
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
func (_m *ConfigProvider) LoadDriverInstance(driverInstanceID string) (*config.DriverInstance, error) {
	ret := _m.Called(driverInstanceID)

	var r0 *config.DriverInstance
	if rf, ok := ret.Get(0).(func(string) *config.DriverInstance); ok {
		r0 = rf(driverInstanceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.DriverInstance)
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
func (_m *ConfigProvider) GetUaaAuthConfig() (*config.UaaAuth, error) {
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
func (_m *ConfigProvider) SetDriver(driverid string, driver config.Driver) error {
	ret := _m.Called(driverid, driver)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, config.Driver) error); ok {
		r0 = rf(driverid, driver)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ConfigProvider) GetDriver(driverid string) (*config.Driver, error) {
	ret := _m.Called(driverid)

	var r0 *config.Driver
	if rf, ok := ret.Get(0).(func(string) *config.Driver); ok {
		r0 = rf(driverid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.Driver)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(driverid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ConfigProvider) DeleteDriver(driverid string) error {
	ret := _m.Called(driverid)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(driverid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ConfigProvider) SetDriverInstance(driverid string, instanceid string, driverInstance config.DriverInstance) error {
	ret := _m.Called(driverid, instanceid, driverInstance)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, config.DriverInstance) error); ok {
		r0 = rf(driverid, instanceid, driverInstance)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ConfigProvider) GetDriverInstance(instanceid string) (*config.DriverInstance, error) {
	ret := _m.Called(instanceid)

	var r0 *config.DriverInstance
	if rf, ok := ret.Get(0).(func(string) *config.DriverInstance); ok {
		r0 = rf(instanceid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.DriverInstance)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(instanceid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ConfigProvider) DeleteDriverInstance(instanceid string) error {
	ret := _m.Called(instanceid)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(instanceid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ConfigProvider) SetService(instanceid string, service brokerapi.Service) error {
	ret := _m.Called(instanceid, service)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, brokerapi.Service) error); ok {
		r0 = rf(instanceid, service)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ConfigProvider) GetService(serviceid string) (*brokerapi.Service, string, error) {
	ret := _m.Called(serviceid)

	var r0 *brokerapi.Service
	if rf, ok := ret.Get(0).(func(string) *brokerapi.Service); ok {
		r0 = rf(serviceid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*brokerapi.Service)
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
func (_m *ConfigProvider) DeleteService(instanceid string) error {
	ret := _m.Called(instanceid)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(instanceid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ConfigProvider) SetDial(instanceid string, dialid string, dial config.Dial) error {
	ret := _m.Called(instanceid, dialid, dial)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, config.Dial) error); ok {
		r0 = rf(instanceid, dialid, dial)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ConfigProvider) GetDial(dialid string) (*config.Dial, error) {
	ret := _m.Called(dialid)

	var r0 *config.Dial
	if rf, ok := ret.Get(0).(func(string) *config.Dial); ok {
		r0 = rf(dialid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.Dial)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(dialid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ConfigProvider) DeleteDial(dialid string) error {
	ret := _m.Called(dialid)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(dialid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ConfigProvider) DriverInstanceNameExists(driverInstanceName string) (bool, error) {
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
func (_m *ConfigProvider) DriverTypeExists(driverType string) (bool, error) {
	ret := _m.Called(driverType)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(driverType)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(driverType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ConfigProvider) GetPlan(plandid string) (*brokerapi.ServicePlan, string, string, error) {
	ret := _m.Called(plandid)

	var r0 *brokerapi.ServicePlan
	if rf, ok := ret.Get(0).(func(string) *brokerapi.ServicePlan); ok {
		r0 = rf(plandid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*brokerapi.ServicePlan)
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
func (_m *ConfigProvider) GetDriversPath() (string, error) {
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
