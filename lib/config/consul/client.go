package consul

import (
	"github.com/hashicorp/consul/api"
	"strings"
)

type ConsulProvisioner struct {
	ConsulClient *api.Client
}

func NewDefault() (ConsulProvisionerInterface, error) {
	var err error
	provisioner := ConsulProvisioner{}

	provisioner.ConsulClient, err = api.NewClient(api.DefaultConfig())

	if err != nil {
		return nil, err
	}

	return &provisioner, nil
}

func New(config *api.Config) (ConsulProvisionerInterface, error) {
	var err error
	provisioner := ConsulProvisioner{}

	provisioner.ConsulClient, err = api.NewClient(config)

	return &provisioner, err
}

func (e *ConsulProvisioner) AddKV(key string, value []byte, options *api.WriteOptions) error {
	kv := e.ConsulClient.KV()

	p := &api.KVPair{Key: key, Value: value}

	_, err := kv.Put(p, options)

	return err
}

func (e *ConsulProvisioner) GetValue(key string) ([]byte, error) {
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

func (e *ConsulProvisioner) PutKVs(pairs *api.KVPairs, options *api.WriteOptions) error {
	var err error
	for _, kv := range *pairs {
		err = e.AddKV(kv.Key, kv.Value, options)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *ConsulProvisioner) GetAllKeys(prefix string, separator string, options *api.QueryOptions) ([]string, error) {
	kv := e.ConsulClient.KV()
	result, _, err := kv.Keys(prefix, separator, options)

	return result, err
}

func (e *ConsulProvisioner) GetAllKVs(prefix string, options *api.QueryOptions) (api.KVPairs, error) {
	kv := e.ConsulClient.KV()
	result, _, err := kv.List(prefix, options)

	return result, err
}

func (e *ConsulProvisioner) DeleteKVs(prefix string, options *api.WriteOptions) error {
	kv := e.ConsulClient.KV()
	_, err := kv.DeleteTree(prefix, options)

	return err
}

func (e *ConsulProvisioner) DeleteKV(key string, options *api.WriteOptions) error {
	kv := e.ConsulClient.KV()

	_, err := kv.Delete(key, options)

	return err
}
