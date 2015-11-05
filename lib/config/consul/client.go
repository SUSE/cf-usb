package consul

import (
	"github.com/hashicorp/consul/api"
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

	if err != nil {
		return nil, err
	}

	return &provisioner, nil
}

func (e *ConsulProvisioner) AddKV(key string, value []byte, options *api.WriteOptions) error {
	kv := e.ConsulClient.KV()

	p := &api.KVPair{Key: key, Value: value}

	_, err := kv.Put(p, options)
	if err != nil {
		return err
	}
	return nil
}

func (e *ConsulProvisioner) GetValue(key string) ([]byte, error) {
	kv := e.ConsulClient.KV()

	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return nil, err
	}

	return pair.Value, nil
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
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e *ConsulProvisioner) GetAllKVs(prefix string, options *api.QueryOptions) (api.KVPairs, error) {
	kv := e.ConsulClient.KV()
	result, _, err := kv.List(prefix, options)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e *ConsulProvisioner) DelteKVs(prefix string, options *api.WriteOptions) error {
	kv := e.ConsulClient.KV()
	_, err := kv.DeleteTree(prefix, options)
	if err != nil {
		return err
	}
	return nil
}

func (e *ConsulProvisioner) DeleteKV(key string, options *api.WriteOptions) error {
	kv := e.ConsulClient.KV()

	_, err := kv.Delete(key, options)

	if err != nil {
		return err
	}
	return nil
}
