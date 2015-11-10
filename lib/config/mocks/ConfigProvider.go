package mocks

import "github.com/hpcloud/cf-usb/lib/config"
import "github.com/stretchr/testify/mock"

import "github.com/pivotal-cf/brokerapi"

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
func (_m *ConfigProvider) GetDriverInstanceConfig(driverInstanceID string) (*config.DriverInstance, error) {
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
func (_m *ConfigProvider) SetDriver(_a0 config.Driver) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(config.Driver) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ConfigProvider) GetDriver(_a0 string) (config.Driver, error) {
	ret := _m.Called(_a0)

	var r0 config.Driver
	if rf, ok := ret.Get(0).(func(string) config.Driver); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(config.Driver)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ConfigProvider) DeleteDriver(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ConfigProvider) SetDriverInstance(_a0 string, _a1 config.DriverInstance) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, config.DriverInstance) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ConfigProvider) GetDriverInstance(_a0 string) (config.DriverInstance, error) {
	ret := _m.Called(_a0)

	var r0 config.DriverInstance
	if rf, ok := ret.Get(0).(func(string) config.DriverInstance); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(config.DriverInstance)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ConfigProvider) DeleteDriverInstance(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ConfigProvider) SetService(_a0 string, _a1 brokerapi.Service) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, brokerapi.Service) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ConfigProvider) GetService(_a0 string) (brokerapi.Service, error) {
	ret := _m.Called(_a0)

	var r0 brokerapi.Service
	if rf, ok := ret.Get(0).(func(string) brokerapi.Service); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(brokerapi.Service)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ConfigProvider) DeleteService(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ConfigProvider) SetDial(_a0 string, _a1 config.Dial) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, config.Dial) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ConfigProvider) GetDial(_a0 string, _a1 string) (config.Dial, error) {
	ret := _m.Called(_a0, _a1)

	var r0 config.Dial
	if rf, ok := ret.Get(0).(func(string, string) config.Dial); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(config.Dial)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ConfigProvider) DeleteDial(_a0 string, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
