package consul

import "github.com/hashicorp/consul/api"

type ConsulProvisionerInterface interface {
	AddKV(string, []byte, *api.WriteOptions) error
	PutKVs(*api.KVPairs, *api.WriteOptions) error
	GetValue(string) ([]byte, error)
	GetAllKVs(string, *api.QueryOptions) (api.KVPairs, error)
	DeleteKV(string, *api.WriteOptions) error
	DelteKVs(string, *api.WriteOptions) error
	GetAllKeys(string, string, *api.QueryOptions) ([]string, error)
}
