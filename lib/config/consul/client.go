package consul

import (
	"strings"

	"github.com/hashicorp/consul/api"
)

//ProvisionerConsul provides the definition of consul provisioner
type ProvisionerConsul struct {
	ConsulClient *api.Client
}

//NewDefault returns an instance of Consul Provisioner with default values
func NewDefault() (Provisioner, error) {
	var err error
	provisioner := ProvisionerConsul{}

	provisioner.ConsulClient, err = api.NewClient(api.DefaultConfig())

	if err != nil {
		return nil, err
	}

	return &provisioner, nil
}

//New returns a new Consul Provisioner
func New(config *api.Config) (Provisioner, error) {
	var err error
	provisioner := ProvisionerConsul{}

	provisioner.ConsulClient, err = api.NewClient(config)

	return &provisioner, err
}

//AddKV adds a key value pair to consul db
func (e *ProvisionerConsul) AddKV(key string, value []byte, options *api.WriteOptions) error {
	kv := e.ConsulClient.KV()

	p := &api.KVPair{Key: key, Value: value}

	_, err := kv.Put(p, options)

	return err
}

//GetValue gets the value with the specified key from consul
func (e *ProvisionerConsul) GetValue(key string) ([]byte, error) {
	kv := e.ConsulClient.KV()

	keys, err := e.GetAllKeys(strings.Split(key, "/")[0], "", nil)
	if err != nil {
		return nil, err
	}
	exists := false
	for _, keyinfo := range keys {
		if strings.Contains(keyinfo, key) {
			exists = true
			break
		}
	}
	if exists == false {
		return nil, nil
	}
	var options api.QueryOptions
	options.AllowStale = false
	pair, _, err := kv.Get(key, &options)

	return pair.Value, err
}

//PutKVs puts several key-value pairs in Consul
func (e *ProvisionerConsul) PutKVs(pairs *api.KVPairs, options *api.WriteOptions) error {
	var err error
	for _, kv := range *pairs {
		err = e.AddKV(kv.Key, kv.Value, options)
		if err != nil {
			return err
		}
	}
	return nil
}

//GetAllKeys obtains all keys from consul
func (e *ProvisionerConsul) GetAllKeys(prefix string, separator string, options *api.QueryOptions) ([]string, error) {
	kv := e.ConsulClient.KV()
	result, _, err := kv.Keys(prefix, separator, options)

	return result, err
}

//GetAllKVs obtains all key-values pairs from consul
func (e *ProvisionerConsul) GetAllKVs(prefix string, options *api.QueryOptions) (api.KVPairs, error) {
	kv := e.ConsulClient.KV()
	result, _, err := kv.List(prefix, options)

	return result, err
}

//DeleteKVs deletes the specified key-value pairs from consul
func (e *ProvisionerConsul) DeleteKVs(prefix string, options *api.WriteOptions) error {
	kv := e.ConsulClient.KV()
	_, err := kv.DeleteTree(prefix, options)

	return err
}

//DeleteKV deletes the specified key-value pair from consul
func (e *ProvisionerConsul) DeleteKV(key string, options *api.WriteOptions) error {
	kv := e.ConsulClient.KV()

	_, err := kv.Delete(key, options)

	return err
}
