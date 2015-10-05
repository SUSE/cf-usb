# boltdb

This package is a gocfbroker.Storer implementation in order to provide storage
backends for Cloud Foundry service brokers.

### Configuration

`boltdb.New()` takes two parameters:

*filename:* Path to the boltdb database filename, can be relative.

*bucket:* The boltdb bucket under which all the keys will reside.
