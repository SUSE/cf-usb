package config

import (
	"github.com/stretchr/testify/assert"
	"testing"

	redisMock "github.com/hpcloud/cf-usb/lib/config/redis/mocks"
	"github.com/pivotal-cf/brokerapi"
	"github.com/stretchr/testify/mock"
)

var RedisTestConfig = struct {
	Provider ConfigProvider
}{}

func Test_Redis_LoadConfiguration(t *testing.T) {
	assert := assert.New(t)
	provisioner := new(redisMock.RedisProvisionerInterface)

	provisioner.On("GetValue", "broker_api").Return("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}", nil)

	provisioner.On("GetValue", "management_api").Return("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}", nil)

	provisioner.On("KeyExists", "drivers").Return(false, nil)

	RedisTestConfig.Provider = NewRedisConfig(provisioner)
	config, err := RedisTestConfig.Provider.LoadConfiguration()
	t.Log(config.BrokerAPI)
	assert.NoError(err)
}

func Test_Redis_GetDriver(t *testing.T) {
	assert := assert.New(t)
	provisioner := new(redisMock.RedisProvisionerInterface)
	provisioner.On("GetValue", "broker_api").Return("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}", nil)

	provisioner.On("GetValue", "management_api").Return("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}", nil)

	provisioner.On("KeyExists", "drivers").Return(true, nil)
	provisioner.On("GetValue", "drivers").Return("{\"00000000-0000-0000-0000-000000000001\":{\"driver_type\":\"dummy\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000002\":{\"name\":\"dummy1\",\"id\":\"A0000000-0000-0000-0000-000000000002\",\"configuration\":{\"property_one\":\"one\",\"property_two\":\"two\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"configuration\":{\"max_dbsize_mb\":2},\"plan\":{\"name\":\"free\",\"id\":\"53425178-F731-49E7-9E53-5CF4BE9D807A\",\"description\":\"This is the first plan\",\"free\":true}},\"B0000000-0000-0000-0000-000000000002\":{\"configuration\":{\"max_dbsize_mb\":100},\"plan\":{\"name\":\"secondary\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F07\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442A\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}},\"A0000000-0000-0000-0000-000000000003\":{\"name\":\"dummy2\",\"configuration\":{\"property_one\":\"onenew\",\"property_two\":\"twonew\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"plandummy2\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442B\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}}}},\"00000000-0000-0000-0000-000000000002\":{\"driver_type\":\"mssql\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000003\":{\"driver_id\":\"00000000-0000-0000-0000-000000000002\",\"name\":\"local-mssql\",\"configuration\":{\"brokerGoSqlDriver\":\"mssql\",\"brokerMssqlConnection\":{\"server\":\"127.0.0.1\",\"port\":\"38017\",\"database\":\"master\",\"user id\":\"sa\",\"password\":\"password1234!\"},\"servedMssqlBindingHostname\":\"192.168.1.10\",\"servedMssqlBindingPort\":38017},\"dials\":{\"C0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"planmssql\",\"id\":\"888B59E0-C2A2-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442C\",\"bindable\":true,\"name\":\"mssql\",\"description\":\"MSSQL Service\",\"tags\":[\"mssql\",\"mssql\"],\"metadata\":{\"providerDisplayName\":\"MSSQL Service Ltd.\"}}}}}}", nil)
	RedisTestConfig.Provider = NewRedisConfig(provisioner)
	driver, err := RedisTestConfig.Provider.GetDriver("00000000-0000-0000-0000-000000000001")
	if driver != nil {
		t.Log((*driver).DriverType)
	}
	assert.NoError(err)
}

func Test_Redis_GetDriverInstance(t *testing.T) {
	assert := assert.New(t)
	provisioner := new(redisMock.RedisProvisionerInterface)
	provisioner.On("GetValue", "broker_api").Return("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}", nil)

	provisioner.On("GetValue", "management_api").Return("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}", nil)

	provisioner.On("KeyExists", "drivers").Return(true, nil)
	provisioner.On("GetValue", "drivers").Return("{\"00000000-0000-0000-0000-000000000001\":{\"driver_type\":\"dummy\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000002\":{\"name\":\"dummy1\",\"id\":\"A0000000-0000-0000-0000-000000000002\",\"configuration\":{\"property_one\":\"one\",\"property_two\":\"two\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"configuration\":{\"max_dbsize_mb\":2},\"plan\":{\"name\":\"free\",\"id\":\"53425178-F731-49E7-9E53-5CF4BE9D807A\",\"description\":\"This is the first plan\",\"free\":true}},\"B0000000-0000-0000-0000-000000000002\":{\"configuration\":{\"max_dbsize_mb\":100},\"plan\":{\"name\":\"secondary\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F07\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442A\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}},\"A0000000-0000-0000-0000-000000000003\":{\"name\":\"dummy2\",\"configuration\":{\"property_one\":\"onenew\",\"property_two\":\"twonew\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"plandummy2\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442B\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}}}},\"00000000-0000-0000-0000-000000000002\":{\"driver_type\":\"mssql\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000003\":{\"driver_id\":\"00000000-0000-0000-0000-000000000002\",\"name\":\"local-mssql\",\"configuration\":{\"brokerGoSqlDriver\":\"mssql\",\"brokerMssqlConnection\":{\"server\":\"127.0.0.1\",\"port\":\"38017\",\"database\":\"master\",\"user id\":\"sa\",\"password\":\"password1234!\"},\"servedMssqlBindingHostname\":\"192.168.1.10\",\"servedMssqlBindingPort\":38017},\"dials\":{\"C0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"planmssql\",\"id\":\"888B59E0-C2A2-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442C\",\"bindable\":true,\"name\":\"mssql\",\"description\":\"MSSQL Service\",\"tags\":[\"mssql\",\"mssql\"],\"metadata\":{\"providerDisplayName\":\"MSSQL Service Ltd.\"}}}}}}", nil)
	RedisTestConfig.Provider = NewRedisConfig(provisioner)
	instance, err := RedisTestConfig.Provider.GetDriverInstance("A0000000-0000-0000-0000-000000000002")
	t.Log((*instance).Name)
	assert.NoError(err)
}

func Test_Redis_GetDial(t *testing.T) {
	assert := assert.New(t)
	provisioner := new(redisMock.RedisProvisionerInterface)
	provisioner.On("GetValue", "broker_api").Return("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}", nil)

	provisioner.On("GetValue", "management_api").Return("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}", nil)

	provisioner.On("KeyExists", "drivers").Return(true, nil)
	provisioner.On("GetValue", "drivers").Return("{\"00000000-0000-0000-0000-000000000001\":{\"driver_type\":\"dummy\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000002\":{\"name\":\"dummy1\",\"id\":\"A0000000-0000-0000-0000-000000000002\",\"configuration\":{\"property_one\":\"one\",\"property_two\":\"two\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"configuration\":{\"max_dbsize_mb\":2},\"plan\":{\"name\":\"free\",\"id\":\"53425178-F731-49E7-9E53-5CF4BE9D807A\",\"description\":\"This is the first plan\",\"free\":true}},\"B0000000-0000-0000-0000-000000000002\":{\"configuration\":{\"max_dbsize_mb\":100},\"plan\":{\"name\":\"secondary\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F07\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442A\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}},\"A0000000-0000-0000-0000-000000000003\":{\"name\":\"dummy2\",\"configuration\":{\"property_one\":\"onenew\",\"property_two\":\"twonew\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"plandummy2\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442B\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}}}},\"00000000-0000-0000-0000-000000000002\":{\"driver_type\":\"mssql\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000003\":{\"driver_id\":\"00000000-0000-0000-0000-000000000002\",\"name\":\"local-mssql\",\"configuration\":{\"brokerGoSqlDriver\":\"mssql\",\"brokerMssqlConnection\":{\"server\":\"127.0.0.1\",\"port\":\"38017\",\"database\":\"master\",\"user id\":\"sa\",\"password\":\"password1234!\"},\"servedMssqlBindingHostname\":\"192.168.1.10\",\"servedMssqlBindingPort\":38017},\"dials\":{\"C0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"planmssql\",\"id\":\"888B59E0-C2A2-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442C\",\"bindable\":true,\"name\":\"mssql\",\"description\":\"MSSQL Service\",\"tags\":[\"mssql\",\"mssql\"],\"metadata\":{\"providerDisplayName\":\"MSSQL Service Ltd.\"}}}}}}", nil)
	RedisTestConfig.Provider = NewRedisConfig(provisioner)
	dial, err := RedisTestConfig.Provider.GetDial("A0000000-0000-0000-0000-000000000002", "B0000000-0000-0000-0000-000000000001")
	if dial != nil {
		t.Log(*dial)
	}
	assert.NoError(err)
}

func Test_Redis_GetService(t *testing.T) {
	assert := assert.New(t)
	provisioner := new(redisMock.RedisProvisionerInterface)
	provisioner.On("GetValue", "broker_api").Return("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}", nil)

	provisioner.On("GetValue", "management_api").Return("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}", nil)

	provisioner.On("KeyExists", "drivers").Return(true, nil)
	provisioner.On("GetValue", "drivers").Return("{\"00000000-0000-0000-0000-000000000001\":{\"driver_type\":\"dummy\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000002\":{\"name\":\"dummy1\",\"id\":\"A0000000-0000-0000-0000-000000000002\",\"configuration\":{\"property_one\":\"one\",\"property_two\":\"two\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"configuration\":{\"max_dbsize_mb\":2},\"plan\":{\"name\":\"free\",\"id\":\"53425178-F731-49E7-9E53-5CF4BE9D807A\",\"description\":\"This is the first plan\",\"free\":true}},\"B0000000-0000-0000-0000-000000000002\":{\"configuration\":{\"max_dbsize_mb\":100},\"plan\":{\"name\":\"secondary\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F07\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442A\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}},\"A0000000-0000-0000-0000-000000000003\":{\"name\":\"dummy2\",\"configuration\":{\"property_one\":\"onenew\",\"property_two\":\"twonew\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"plandummy2\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442B\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}}}},\"00000000-0000-0000-0000-000000000002\":{\"driver_type\":\"mssql\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000003\":{\"driver_id\":\"00000000-0000-0000-0000-000000000002\",\"name\":\"local-mssql\",\"configuration\":{\"brokerGoSqlDriver\":\"mssql\",\"brokerMssqlConnection\":{\"server\":\"127.0.0.1\",\"port\":\"38017\",\"database\":\"master\",\"user id\":\"sa\",\"password\":\"password1234!\"},\"servedMssqlBindingHostname\":\"192.168.1.10\",\"servedMssqlBindingPort\":38017},\"dials\":{\"C0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"planmssql\",\"id\":\"888B59E0-C2A2-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442C\",\"bindable\":true,\"name\":\"mssql\",\"description\":\"MSSQL Service\",\"tags\":[\"mssql\",\"mssql\"],\"metadata\":{\"providerDisplayName\":\"MSSQL Service Ltd.\"}}}}}}", nil)

	RedisTestConfig.Provider = NewRedisConfig(provisioner)
	service, err := RedisTestConfig.Provider.GetService("A0000000-0000-0000-0000-000000000002")
	if service != nil {
		t.Log(*service)
	}
	assert.NoError(err)
}

func Test_Redis_SetDriver(t *testing.T) {
	assert := assert.New(t)
	provisioner := new(redisMock.RedisProvisionerInterface)

	provisioner.On("GetValue", "broker_api").Return("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}", nil)

	provisioner.On("GetValue", "management_api").Return("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}", nil)

	provisioner.On("KeyExists", "drivers").Return(false, nil)

	provisioner.On("SetKV", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	RedisTestConfig.Provider = NewRedisConfig(provisioner)

	var driver Driver
	driver.DriverType = "testDriver"

	err := RedisTestConfig.Provider.SetDriver("testDriverID", driver)
	assert.NoError(err)
}

/*
func Test_Redis_SetDriverInstance(t *testing.T) {
	assert := assert.New(t)
	provisioner := new(redisMock.RedisProvisionerInterface)

	provisioner.On("GetValue", "broker_api").Return("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}", nil)

	provisioner.On("GetValue", "management_api").Return("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}", nil)

	provisioner.On("KeyExists", "drivers").Return(true, nil)
	provisioner.On("GetValue", "drivers").Return("{\"00000000-0000-0000-0000-000000000001\":{\"driver_type\":\"dummy\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000002\":{\"name\":\"dummy1\",\"id\":\"A0000000-0000-0000-0000-000000000002\",\"configuration\":{\"property_one\":\"one\",\"property_two\":\"two\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"configuration\":{\"max_dbsize_mb\":2},\"plan\":{\"name\":\"free\",\"id\":\"53425178-F731-49E7-9E53-5CF4BE9D807A\",\"description\":\"This is the first plan\",\"free\":true}},\"B0000000-0000-0000-0000-000000000002\":{\"configuration\":{\"max_dbsize_mb\":100},\"plan\":{\"name\":\"secondary\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F07\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442A\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}},\"A0000000-0000-0000-0000-000000000003\":{\"name\":\"dummy2\",\"configuration\":{\"property_one\":\"onenew\",\"property_two\":\"twonew\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"plandummy2\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442B\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}}}},\"00000000-0000-0000-0000-000000000002\":{\"driver_type\":\"mssql\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000003\":{\"driver_id\":\"00000000-0000-0000-0000-000000000002\",\"name\":\"local-mssql\",\"configuration\":{\"brokerGoSqlDriver\":\"mssql\",\"brokerMssqlConnection\":{\"server\":\"127.0.0.1\",\"port\":\"38017\",\"database\":\"master\",\"user id\":\"sa\",\"password\":\"password1234!\"},\"servedMssqlBindingHostname\":\"192.168.1.10\",\"servedMssqlBindingPort\":38017},\"dials\":{\"C0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"planmssql\",\"id\":\"888B59E0-C2A2-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442C\",\"bindable\":true,\"name\":\"mssql\",\"description\":\"MSSQL Service\",\"tags\":[\"mssql\",\"mssql\"],\"metadata\":{\"providerDisplayName\":\"MSSQL Service Ltd.\"}}}}}}", nil)

	provisioner.On("SetKV", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	RedisTestConfig.Provider = NewRedisConfig(provisioner)

	var instance DriverInstance
	instance.ID = "I0000000-0000-0000-0000-0000000000T1"
	instance.Name = "testInstance"

	err := RedisTestConfig.Provider.SetDriverInstance("00000000-0000-0000-0000-000000000001", instance)
	assert.NoError(err)
}*/

func Test_Redis_SetDial(t *testing.T) {
	assert := assert.New(t)
	provisioner := new(redisMock.RedisProvisionerInterface)

	provisioner.On("GetValue", "broker_api").Return("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}", nil)

	provisioner.On("GetValue", "management_api").Return("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}", nil)

	provisioner.On("KeyExists", "drivers").Return(true, nil)
	provisioner.On("GetValue", "drivers").Return("{\"00000000-0000-0000-0000-000000000001\":{\"driver_type\":\"dummy\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000002\":{\"name\":\"dummy1\",\"id\":\"A0000000-0000-0000-0000-000000000002\",\"configuration\":{\"property_one\":\"one\",\"property_two\":\"two\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"configuration\":{\"max_dbsize_mb\":2},\"plan\":{\"name\":\"free\",\"id\":\"53425178-F731-49E7-9E53-5CF4BE9D807A\",\"description\":\"This is the first plan\",\"free\":true}},\"B0000000-0000-0000-0000-000000000002\":{\"configuration\":{\"max_dbsize_mb\":100},\"plan\":{\"name\":\"secondary\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F07\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442A\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}},\"A0000000-0000-0000-0000-000000000003\":{\"name\":\"dummy2\",\"configuration\":{\"property_one\":\"onenew\",\"property_two\":\"twonew\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"plandummy2\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442B\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}}}},\"00000000-0000-0000-0000-000000000002\":{\"driver_type\":\"mssql\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000003\":{\"driver_id\":\"00000000-0000-0000-0000-000000000002\",\"name\":\"local-mssql\",\"configuration\":{\"brokerGoSqlDriver\":\"mssql\",\"brokerMssqlConnection\":{\"server\":\"127.0.0.1\",\"port\":\"38017\",\"database\":\"master\",\"user id\":\"sa\",\"password\":\"password1234!\"},\"servedMssqlBindingHostname\":\"192.168.1.10\",\"servedMssqlBindingPort\":38017},\"dials\":{\"C0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"planmssql\",\"id\":\"888B59E0-C2A2-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442C\",\"bindable\":true,\"name\":\"mssql\",\"description\":\"MSSQL Service\",\"tags\":[\"mssql\",\"mssql\"],\"metadata\":{\"providerDisplayName\":\"MSSQL Service Ltd.\"}}}}}}", nil)

	provisioner.On("SetKV", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	RedisTestConfig.Provider = NewRedisConfig(provisioner)

	var dial Dial

	err := RedisTestConfig.Provider.SetDial("A0000000-0000-0000-0000-000000000002", "testDialID", dial)
	assert.NoError(err)
}

func Test_Redis_SetService(t *testing.T) {
	assert := assert.New(t)
	provisioner := new(redisMock.RedisProvisionerInterface)

	provisioner.On("GetValue", "broker_api").Return("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}", nil)

	provisioner.On("GetValue", "management_api").Return("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}", nil)

	provisioner.On("KeyExists", "drivers").Return(true, nil)
	provisioner.On("GetValue", "drivers").Return("{\"00000000-0000-0000-0000-000000000001\":{\"driver_type\":\"dummy\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000002\":{\"name\":\"dummy1\",\"id\":\"A0000000-0000-0000-0000-000000000002\",\"configuration\":{\"property_one\":\"one\",\"property_two\":\"two\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"configuration\":{\"max_dbsize_mb\":2},\"plan\":{\"name\":\"free\",\"id\":\"53425178-F731-49E7-9E53-5CF4BE9D807A\",\"description\":\"This is the first plan\",\"free\":true}},\"B0000000-0000-0000-0000-000000000002\":{\"configuration\":{\"max_dbsize_mb\":100},\"plan\":{\"name\":\"secondary\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F07\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442A\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}},\"A0000000-0000-0000-0000-000000000003\":{\"name\":\"dummy2\",\"configuration\":{\"property_one\":\"onenew\",\"property_two\":\"twonew\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"plandummy2\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442B\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}}}},\"00000000-0000-0000-0000-000000000002\":{\"driver_type\":\"mssql\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000003\":{\"driver_id\":\"00000000-0000-0000-0000-000000000002\",\"name\":\"local-mssql\",\"configuration\":{\"brokerGoSqlDriver\":\"mssql\",\"brokerMssqlConnection\":{\"server\":\"127.0.0.1\",\"port\":\"38017\",\"database\":\"master\",\"user id\":\"sa\",\"password\":\"password1234!\"},\"servedMssqlBindingHostname\":\"192.168.1.10\",\"servedMssqlBindingPort\":38017},\"dials\":{\"C0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"planmssql\",\"id\":\"888B59E0-C2A2-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442C\",\"bindable\":true,\"name\":\"mssql\",\"description\":\"MSSQL Service\",\"tags\":[\"mssql\",\"mssql\"],\"metadata\":{\"providerDisplayName\":\"MSSQL Service Ltd.\"}}}}}}", nil)

	provisioner.On("SetKV", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	RedisTestConfig.Provider = NewRedisConfig(provisioner)

	var service brokerapi.Service
	service.Bindable = true
	service.ID = "testServiceID"
	service.Description = "test service"
	service.Name = "testService2"
	service.Tags = []string{"test"}
	err := RedisTestConfig.Provider.SetService("A0000000-0000-0000-0000-000000000002", service)
	assert.NoError(err)
}

func Test_Redis_DeleteDriver(t *testing.T) {
	assert := assert.New(t)
	provisioner := new(redisMock.RedisProvisionerInterface)

	provisioner.On("GetValue", "broker_api").Return("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}", nil)

	provisioner.On("GetValue", "management_api").Return("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}", nil)

	provisioner.On("KeyExists", "drivers").Return(true, nil)
	provisioner.On("GetValue", "drivers").Return("{\"00000000-0000-0000-0000-000000000001\":{\"driver_type\":\"dummy\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000002\":{\"name\":\"dummy1\",\"id\":\"A0000000-0000-0000-0000-000000000002\",\"configuration\":{\"property_one\":\"one\",\"property_two\":\"two\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"configuration\":{\"max_dbsize_mb\":2},\"plan\":{\"name\":\"free\",\"id\":\"53425178-F731-49E7-9E53-5CF4BE9D807A\",\"description\":\"This is the first plan\",\"free\":true}},\"B0000000-0000-0000-0000-000000000002\":{\"configuration\":{\"max_dbsize_mb\":100},\"plan\":{\"name\":\"secondary\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F07\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442A\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}},\"A0000000-0000-0000-0000-000000000003\":{\"name\":\"dummy2\",\"configuration\":{\"property_one\":\"onenew\",\"property_two\":\"twonew\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"plandummy2\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442B\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}}}},\"00000000-0000-0000-0000-000000000002\":{\"driver_type\":\"mssql\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000003\":{\"driver_id\":\"00000000-0000-0000-0000-000000000002\",\"name\":\"local-mssql\",\"configuration\":{\"brokerGoSqlDriver\":\"mssql\",\"brokerMssqlConnection\":{\"server\":\"127.0.0.1\",\"port\":\"38017\",\"database\":\"master\",\"user id\":\"sa\",\"password\":\"password1234!\"},\"servedMssqlBindingHostname\":\"192.168.1.10\",\"servedMssqlBindingPort\":38017},\"dials\":{\"C0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"planmssql\",\"id\":\"888B59E0-C2A2-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442C\",\"bindable\":true,\"name\":\"mssql\",\"description\":\"MSSQL Service\",\"tags\":[\"mssql\",\"mssql\"],\"metadata\":{\"providerDisplayName\":\"MSSQL Service Ltd.\"}}}}}}", nil)

	provisioner.On("SetKV", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	RedisTestConfig.Provider = NewRedisConfig(provisioner)

	err := RedisTestConfig.Provider.DeleteDriver("00000000-0000-0000-0000-000000000001")
	assert.NoError(err)
}

func Test_Redis_DeleteDriverInstance(t *testing.T) {
	assert := assert.New(t)
	provisioner := new(redisMock.RedisProvisionerInterface)

	provisioner.On("GetValue", "broker_api").Return("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}", nil)

	provisioner.On("GetValue", "management_api").Return("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}", nil)

	provisioner.On("KeyExists", "drivers").Return(true, nil)
	provisioner.On("GetValue", "drivers").Return("{\"00000000-0000-0000-0000-000000000001\":{\"driver_type\":\"dummy\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000002\":{\"name\":\"dummy1\",\"id\":\"A0000000-0000-0000-0000-000000000002\",\"configuration\":{\"property_one\":\"one\",\"property_two\":\"two\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"configuration\":{\"max_dbsize_mb\":2},\"plan\":{\"name\":\"free\",\"id\":\"53425178-F731-49E7-9E53-5CF4BE9D807A\",\"description\":\"This is the first plan\",\"free\":true}},\"B0000000-0000-0000-0000-000000000002\":{\"configuration\":{\"max_dbsize_mb\":100},\"plan\":{\"name\":\"secondary\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F07\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442A\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}},\"A0000000-0000-0000-0000-000000000003\":{\"name\":\"dummy2\",\"configuration\":{\"property_one\":\"onenew\",\"property_two\":\"twonew\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"plandummy2\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442B\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}}}},\"00000000-0000-0000-0000-000000000002\":{\"driver_type\":\"mssql\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000003\":{\"driver_id\":\"00000000-0000-0000-0000-000000000002\",\"name\":\"local-mssql\",\"configuration\":{\"brokerGoSqlDriver\":\"mssql\",\"brokerMssqlConnection\":{\"server\":\"127.0.0.1\",\"port\":\"38017\",\"database\":\"master\",\"user id\":\"sa\",\"password\":\"password1234!\"},\"servedMssqlBindingHostname\":\"192.168.1.10\",\"servedMssqlBindingPort\":38017},\"dials\":{\"C0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"planmssql\",\"id\":\"888B59E0-C2A2-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442C\",\"bindable\":true,\"name\":\"mssql\",\"description\":\"MSSQL Service\",\"tags\":[\"mssql\",\"mssql\"],\"metadata\":{\"providerDisplayName\":\"MSSQL Service Ltd.\"}}}}}}", nil)

	provisioner.On("SetKV", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	RedisTestConfig.Provider = NewRedisConfig(provisioner)

	err := RedisTestConfig.Provider.DeleteDriverInstance("A0000000-0000-0000-0000-000000000002")
	assert.NoError(err)
}

func Test_Redis_DeleteDial(t *testing.T) {
	assert := assert.New(t)
	provisioner := new(redisMock.RedisProvisionerInterface)

	provisioner.On("GetValue", "broker_api").Return("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}", nil)

	provisioner.On("GetValue", "management_api").Return("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}", nil)

	provisioner.On("KeyExists", "drivers").Return(true, nil)
	provisioner.On("GetValue", "drivers").Return("{\"00000000-0000-0000-0000-000000000001\":{\"driver_type\":\"dummy\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000002\":{\"name\":\"dummy1\",\"id\":\"A0000000-0000-0000-0000-000000000002\",\"configuration\":{\"property_one\":\"one\",\"property_two\":\"two\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"configuration\":{\"max_dbsize_mb\":2},\"plan\":{\"name\":\"free\",\"id\":\"53425178-F731-49E7-9E53-5CF4BE9D807A\",\"description\":\"This is the first plan\",\"free\":true}},\"B0000000-0000-0000-0000-000000000002\":{\"configuration\":{\"max_dbsize_mb\":100},\"plan\":{\"name\":\"secondary\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F07\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442A\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}},\"A0000000-0000-0000-0000-000000000003\":{\"name\":\"dummy2\",\"configuration\":{\"property_one\":\"onenew\",\"property_two\":\"twonew\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"plandummy2\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442B\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}}}},\"00000000-0000-0000-0000-000000000002\":{\"driver_type\":\"mssql\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000003\":{\"driver_id\":\"00000000-0000-0000-0000-000000000002\",\"name\":\"local-mssql\",\"configuration\":{\"brokerGoSqlDriver\":\"mssql\",\"brokerMssqlConnection\":{\"server\":\"127.0.0.1\",\"port\":\"38017\",\"database\":\"master\",\"user id\":\"sa\",\"password\":\"password1234!\"},\"servedMssqlBindingHostname\":\"192.168.1.10\",\"servedMssqlBindingPort\":38017},\"dials\":{\"C0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"planmssql\",\"id\":\"888B59E0-C2A2-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442C\",\"bindable\":true,\"name\":\"mssql\",\"description\":\"MSSQL Service\",\"tags\":[\"mssql\",\"mssql\"],\"metadata\":{\"providerDisplayName\":\"MSSQL Service Ltd.\"}}}}}}", nil)

	provisioner.On("SetKV", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	RedisTestConfig.Provider = NewRedisConfig(provisioner)

	err := RedisTestConfig.Provider.DeleteDial("A0000000-0000-0000-0000-000000000002", "B0000000-0000-0000-0000-000000000001")
	assert.NoError(err)
}

func Test_Redis_DeleteService(t *testing.T) {
	assert := assert.New(t)
	provisioner := new(redisMock.RedisProvisionerInterface)

	provisioner.On("GetValue", "broker_api").Return("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}", nil)

	provisioner.On("GetValue", "management_api").Return("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}", nil)

	provisioner.On("KeyExists", "drivers").Return(true, nil)
	provisioner.On("GetValue", "drivers").Return("{\"00000000-0000-0000-0000-000000000001\":{\"driver_type\":\"dummy\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000002\":{\"name\":\"dummy1\",\"id\":\"A0000000-0000-0000-0000-000000000002\",\"configuration\":{\"property_one\":\"one\",\"property_two\":\"two\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"configuration\":{\"max_dbsize_mb\":2},\"plan\":{\"name\":\"free\",\"id\":\"53425178-F731-49E7-9E53-5CF4BE9D807A\",\"description\":\"This is the first plan\",\"free\":true}},\"B0000000-0000-0000-0000-000000000002\":{\"configuration\":{\"max_dbsize_mb\":100},\"plan\":{\"name\":\"secondary\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F07\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442A\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}},\"A0000000-0000-0000-0000-000000000003\":{\"name\":\"dummy2\",\"configuration\":{\"property_one\":\"onenew\",\"property_two\":\"twonew\"},\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"plandummy2\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442B\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}}}},\"00000000-0000-0000-0000-000000000002\":{\"driver_type\":\"mssql\",\"driver_instances\":{\"A0000000-0000-0000-0000-000000000003\":{\"driver_id\":\"00000000-0000-0000-0000-000000000002\",\"name\":\"local-mssql\",\"configuration\":{\"brokerGoSqlDriver\":\"mssql\",\"brokerMssqlConnection\":{\"server\":\"127.0.0.1\",\"port\":\"38017\",\"database\":\"master\",\"user id\":\"sa\",\"password\":\"password1234!\"},\"servedMssqlBindingHostname\":\"192.168.1.10\",\"servedMssqlBindingPort\":38017},\"dials\":{\"C0000000-0000-0000-0000-000000000001\":{\"plan\":{\"name\":\"planmssql\",\"id\":\"888B59E0-C2A2-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442C\",\"bindable\":true,\"name\":\"mssql\",\"description\":\"MSSQL Service\",\"tags\":[\"mssql\",\"mssql\"],\"metadata\":{\"providerDisplayName\":\"MSSQL Service Ltd.\"}}}}}}", nil)

	provisioner.On("SetKV", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	RedisTestConfig.Provider = NewRedisConfig(provisioner)

	err := RedisTestConfig.Provider.DeleteService("A0000000-0000-0000-0000-000000000002")
	assert.NoError(err)
}
