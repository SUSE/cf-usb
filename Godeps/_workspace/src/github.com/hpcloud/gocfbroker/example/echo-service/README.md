# Cloud Foundry v2 Echo Service Broker

This is a Cloud Foundry v2 Service Broker that uses the service broker library
found at: `github.com/hpcloud/gocfbroker`

It implements an "echo service" which doesn't provision anything but simply
logs to stdout the operations the cloud controller is requesting of it.

## Install

**Requirements:**

* Go 1.4
* Godep

Build uses godep to manage dependencies. If you don't have it get it first:

```bash
go get github.com/tools/godep
```

Then to build the echo service broker:

```bash
# If you haven't cloned the repo do this first
mkdir -p $GOPATH/src/github.com/hpcloud
git clone git@github.com:hpcloud/gocfbroker.git $GOPATH/src/github.com/hpcloud/gocfbroker

# Have to CD here for godep
cd $GOPATH/src/github.com/hpcloud/gocfbroker/example/echo-service
godep go install
```

## Run

```bash
# Ensure config.json is in CWD, see Configuration section below
echo-service # If $GOPATH/bin is not in your path, prefix this command by that.
```

## Configuration

The configuration file is located at [./config.json](config.json)
and contains all the values necessary to run an example service broker.

The listen configuration value is "go style" listen, so the following are valid:

* localhost:5780
* :5780
* 192.168.1.100:5780
* [2001:db8::ff00:42:8329]:5780

**Note:** The db_encryption_key value should be changed to something
more secure (randomly generated). The key length should be 16, 24, or 32 bytes.
This package uses the Go crypto/aes package, see those docs for more information.
