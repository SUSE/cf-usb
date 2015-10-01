package config

import (
	"log"

	"gopkg.in/redis.v2"
)

type redisConfig struct {
	ConfigProvider,
	redisOptions redis.Options
}

func NewRedisConfig(redisOptions redis.Options) ConfigProvider {
	return &redisConfig{redisOptions: redisOptions}
}

func (c *redisConfig) LoadConfiguration() (Config, error) {
	var config Config

	log.Println("Not implemented")

	return config, nil
}

func (c *redisConfig) GetDriverProperties(serviceName string) (DriverProperties, error) {
	var driverProperties DriverProperties

	log.Println("Not implemented")

	return driverProperties, nil

}
func (c *redisConfig) GetDriverTypes() ([]string, error) {
	var driverTypes []string
	log.Println("Not implemented")

	return driverTypes, nil
}
