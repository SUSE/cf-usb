package mocks

import "github.com/stretchr/testify/mock"

import "github.com/hashicorp/consul/api"

type ConsulProvisionerInterface struct {
	mock.Mock
}

func (_m *ConsulProvisionerInterface) AddKV(_a0 string, _a1 []byte, _a2 *api.WriteOptions) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte, *api.WriteOptions) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ConsulProvisionerInterface) PutKVs(_a0 *api.KVPairs, _a1 *api.WriteOptions) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*api.KVPairs, *api.WriteOptions) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ConsulProvisionerInterface) GetValue(_a0 string) ([]byte, error) {
	ret := _m.Called(_a0)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ConsulProvisionerInterface) GetAllKVs(_a0 string, _a1 *api.QueryOptions) (api.KVPairs, error) {
	ret := _m.Called(_a0, _a1)

	var r0 api.KVPairs
	if rf, ok := ret.Get(0).(func(string, *api.QueryOptions) api.KVPairs); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(api.KVPairs)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *api.QueryOptions) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ConsulProvisionerInterface) DeleteKV(_a0 string, _a1 *api.WriteOptions) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *api.WriteOptions) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ConsulProvisionerInterface) DelteKVs(_a0 string, _a1 *api.WriteOptions) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *api.WriteOptions) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ConsulProvisionerInterface) GetAllKeys(_a0 string, _a1 string, _a2 *api.QueryOptions) ([]string, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string, string, *api.QueryOptions) []string); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, *api.QueryOptions) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
