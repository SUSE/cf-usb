package redis

import (
	"time"

	"gopkg.in/redis.v3"
)

//ProvisionerRedis provides the definition of redis provisioner
type ProvisionerRedis struct {
	RedisClient *redis.Client
}

//New creates a new redis Provisioner and returns it or an error if it fails
func New(address string, password string, db int64) (Provisioner, error) {

	provisioner := ProvisionerRedis{}

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

//GetValue gets the value corresponding to the passed key from redis
func (e ProvisionerRedis) GetValue(key string) (string, error) {
	value, err := e.RedisClient.Get(key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

//SetKV sets the value corresponding to the passed key in redis
func (e ProvisionerRedis) SetKV(key string, value string, expiration time.Duration) error {
	_, err := e.RedisClient.Set(key, value, expiration).Result()
	return err
}

//KeyExists checks if the passed key exists already in redis
func (e ProvisionerRedis) KeyExists(key string) (bool, error) {
	exists, err := e.RedisClient.Exists(key).Result()
	if err != nil {
		return false, err
	}
	return exists, nil
}

//RemoveKey removes the passed key from redis
func (e ProvisionerRedis) RemoveKey(key string) (bool, error) {
	_, err := e.RedisClient.Del(key).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}
