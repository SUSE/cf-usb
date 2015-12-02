package redis

import "time"

type RedisProvisionerInterface interface {
	SetKV(string, string, time.Duration) error
	GetValue(string) (string, error)
	KeyExists(string) (bool, error)
}
