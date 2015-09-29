package config

import "gopkg.in/redis.v2"

type redisConfig struct {
	ConfigProvider,
	redisOptions redis.Options
}

func NewRedisConfig(redisOptions redis.Options) ConfigProvider {
	return &redisConfig{redisOptions: redisOptions}
}

func (c *redisConfig) LoadConfiguration() (Config, error) {
	var config Config

	panic("Not Implemented")

	return config, nil
}

func (c *redisConfig) GetDriverProperties(serviceName string) (DriverProperties, error) {
	var driverProperties DriverProperties

	panic("Not Implemented")

	return driverProperties, nil

}
func (c *redisConfig) GetDriverTypes() ([]string, error) {
	var driverTypes []string
	panic("Not Implemented")

	return driverTypes, nil
}
