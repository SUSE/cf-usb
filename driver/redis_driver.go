package driver

var redisDriver RedisDriver

type RedisDriver struct {
}

func (driver *RedisDriver) Provision() error {
	return nil

}
func (driver *RedisDriver) Deprovision() error {
	return nil
}
func (driver *RedisDriver) Bind() error {
	return nil
}
func (driver *RedisDriver) Unbind() error {
	return nil
}
func (driver *RedisDriver) Update() error {
	return nil
}
func (driver *RedisDriver) GetCatalog() error {
	return nil
}
func (driver *RedisDriver) GetInstances() error {
	return nil
}
