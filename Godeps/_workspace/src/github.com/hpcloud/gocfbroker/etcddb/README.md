# etcddb

This package is a gocfbroker.Storer implementation in order to provide storage
backends for Cloud Foundry service brokers.

### Configuration

`etcddb.New()` takes two configuration parameters:

*machines:* A list of machines on which etcd is running that the etcd client
can try to connect to.

*directory:* The etcddb directory key that each of this service broker's related keys
should be stored under.
