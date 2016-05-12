package redis

import "time"

//Provisioner is the interface to use for a provisioner based on redis
type Provisioner interface {
	SetKV(string, string, time.Duration) error
	GetValue(string) (string, error)
	KeyExists(string) (bool, error)
	RemoveKey(string) (bool, error)
}
