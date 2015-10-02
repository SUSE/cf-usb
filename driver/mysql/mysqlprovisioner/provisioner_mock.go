package mysqlprovisioner

import "github.com/stretchr/testify/mock"

import "database/sql"

type MysqlProvisionerMock struct {
	mock.Mock
}

func NewMockProvisioner(username string, password string, host string) MysqlProvisionerInterface {
	return &MysqlProvisionerMock{}
}

func (m *MysqlProvisionerMock) CreateDatabase(_a0 string) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
func (m *MysqlProvisionerMock) DeleteDatabase(_a0 string) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
func (m *MysqlProvisionerMock) Query(_a0 string) (*sql.Rows, error) {
	ret := m.Called(_a0)

	var r0 *sql.Rows
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*sql.Rows)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *MysqlProvisionerMock) CreateUser(_a0 string, _a1 string, _a2 string) error {
	ret := m.Called(_a0, _a1, _a2)

	r0 := ret.Error(0)

	return r0
}
func (m *MysqlProvisionerMock) DeleteUser(_a0 string) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
