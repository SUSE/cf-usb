package redis

import (
	"gopkg.in/redis.v3"
	"time"
)

type RedisProvisioner struct {
	RedisClient *redis.Client
}

func New(address string, password string, db int64) (RedisProvisionerInterface, error) {

	provisioner := RedisProvisioner{}

	provisioner.RedisClient = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	_, err := provisioner.RedisClient.Ping().Result()

	if err != nil {
		return nil, err
	}
	return provisioner, nil
}

func (e RedisProvisioner) GetValue(key string) (string, error) {
	value, err := e.RedisClient.Get(key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (e RedisProvisioner) SetKV(key string, value string, expiration time.Duration) error {
	_, err := e.RedisClient.Set(key, value, expiration).Result()
	return err
}

func (e RedisProvisioner) KeyExists(key string) (bool, error) {
	exists, err := e.RedisClient.Exists(key).Result()
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (e RedisProvisioner) RemoveKey(key string) (bool, error) {
	_, err := e.RedisClient.Del(key).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}
