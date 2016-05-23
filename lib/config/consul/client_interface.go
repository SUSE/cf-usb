package consul

import "github.com/hashicorp/consul/api"

//Provisioner defines a type containing the values needed for ConsulProvisioner
type Provisioner interface {
	AddKV(string, []byte, *api.WriteOptions) error
	PutKVs(*api.KVPairs, *api.WriteOptions) error
	GetValue(string) ([]byte, error)
	GetAllKVs(string, *api.QueryOptions) (api.KVPairs, error)
	DeleteKV(string, *api.WriteOptions) error
	DeleteKVs(string, *api.WriteOptions) error
	GetAllKeys(string, string, *api.QueryOptions) ([]string, error)
}
