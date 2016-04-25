package mocks

import "github.com/stretchr/testify/mock"

type MssqlProvisionerInterface struct {
	mock.Mock
}

func (_m *MssqlProvisionerInterface) Connect(goSqlDriver string, connectionParams map[string]string) error {
	ret := _m.Called(goSqlDriver, connectionParams)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, map[string]string) error); ok {
		r0 = rf(goSqlDriver, connectionParams)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *MssqlProvisionerInterface) IsDatabaseCreated(databaseId string) (bool, error) {
	ret := _m.Called(databaseId)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(databaseId)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(databaseId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MssqlProvisionerInterface) IsUserCreated(databaseId string, userId string) (bool, error) {
	ret := _m.Called(databaseId, userId)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(databaseId, userId)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(databaseId, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MssqlProvisionerInterface) CreateDatabase(databaseId string) error {
	ret := _m.Called(databaseId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(databaseId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *MssqlProvisionerInterface) DeleteDatabase(databaseId string) error {
	ret := _m.Called(databaseId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(databaseId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *MssqlProvisionerInterface) CreateUser(databaseId string, userId string, password string) error {
	ret := _m.Called(databaseId, userId, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(databaseId, userId, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *MssqlProvisionerInterface) DeleteUser(databaseId string, userId string) error {
	ret := _m.Called(databaseId, userId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(databaseId, userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *MssqlProvisionerInterface) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
