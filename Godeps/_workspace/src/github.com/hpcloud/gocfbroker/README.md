# Cloud Foundry v2 Service Broker Library

This library provides a v2 broker library to develop any custom service on
top of any Cloud Foundry based system in Go.

## Build

**Requirements:**

* Go 1.4
* Godep

Build uses godep to manage dependencies. If you don't have it get it first:

```bash
go get github.com/tools/godep
```

Then to build the service broker library:

```bash
mkdir -p $GOPATH/src/github.com/hpcloud
git clone git@github.com:hpcloud/gocfbroker.git $GOPATH/src/github.com/hpcloud/gocfbroker

# Have to CD here for godep
cd $GOPATH/src/github.com/hpcloud/gocfbroker
godep get        # Fetch dependencies
godep go install # Install library
```

**Note:** This library has a godeps file but does not vendor the dependencies
to avoid duplication with binaries produced with it. By building this library
properly first, and using Godeps in your own service broker and vendoring the
dependencies things should continue to work properly. See the example below for
more details.

## Example

There is a full example "echo service" demonstrating use of this library in 
`example/echo-service`

## Usage

1. Build the library (above)
2. Import library in your own package (github.com/hpcloud/gocfbroker)
3. Implement the `Provisioner` interface
4. Choose `boltdb` or `etcddb` OR: Implement your own `Storer`/`StoreOpener`
5. Create a broker instance: `gocfbroker.New` supplying provisioner, db, and
   configuration file.
6. Use the broker type's `Start`.

Mentioned above are three important interfaces to be aware of to use this
library:

**Provisioner**:
This interface implementation is required. The provisioner interface allows
the broker library to use the passed in type to create and delete service
instances and bindings for the service.

**Storer**:
The `Storer` provides generic Key/Value store functionality with some transactional
awareness to ensure robustness in the broker's data operations. The broker
stores its state (instances, bindings, running jobs) in this key/value store.
Because it's an interface it can be implemented in any way the user chooses.
Some implementations exist already (see: [boltdb](boltdb) [etcddb](etcddb)) for
some pre-built options. These also serve as documentation for users wanting to
create their own implementation.

## Configuration

Pass in an options struct to `gocfbroker.New()`. The library implementer can
use any mechanism to fill in the configuration, but two helper functions are
provided to read in JSON configurations: `gocfbroker.LoadConfig()` and
`gocfbroker.LoadConfigReader()`.

A typical configuration struct will include `gocfbroker.Options` so that all
configuration can be stored in one place and read/loaded at the same time.

Example Configuration can be found here: `example/echo-service/config.json`

```go
// myConfig is our service's custom config, we need a debug_level for our
// broker.
type myConfig struct {
  BoltFilename string `json:"bolt_filename"`
  BoltBucket   string `json:"bolt_bucket"`
  DebugLevel   string `json:"debug_level"`
  gocfbroker.Options
}

var config myConfig
if err := LoadConfig("config.json", &config); err != nil {
  // handle
}

db, err := boltdb.New(config.BoltFilename, config.BoltBucket)
if err != nil {
  // handle
}

broker, err := gocfbroker.New(echoService, db, config.Options)
```

## License

Apache 2.0, see LICENSE file.
