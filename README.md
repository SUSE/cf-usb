# Universal service broker

## Summary
The cf-usb project implements and exposes the Cloud Foundry [Service Broker API](http://docs.cloudfoundry.org/services/api.html). 

It uses plugins (drivers) to connect to different services.

![cf-usb](https://region-b.geo-1.objects.hpcloudsvc.com/v1/11899734124432/imgs/usb.png)

## Configuration

**cf-usb** works with two configuration prividers:

- **file configuration provider** - cf-usb can take it's configuration from a .json file. The json file format is similar to a standard broker configuration file. To start the broker with a file configuration provider you must provide the following cli options:

```sh
 ./usb fileConfigProvider --path {path_to_jsonfile}
```

- **redis configuration provider** - cf-usb can take it's configuration from a redis database. By default, it tries to take the redis uri from the $USB_REDIS_URI environment variable, this can be overridden by providing the following cli options:

```sh
 ./usb redisConfigProvider --uri {uri_to_redis}
```

The driver specific properties can be specified in the *driver_configs* array element of the configuration:

| Property | Description | Mandatory |
| :---: | :---: | :---: |
| *driver_type* | type of the driver  | yes |
| *configuration* | driver specific configurations, ex: service connection parameters | yes
| *service_ids* | array of services that are using this driver | yes |

## Drivers

The drivers are plugins that are used by USB to do service specific operations. A driver must implement the following interface:
```sh
	Init(config.DriverProperties, *string) error
	Provision(model.DriverProvisionRequest, *string) error
	Deprovision(model.DriverDeprovisionRequest, *string) error
	Update(model.DriverUpdateRequest, *string) error
	Bind(model.DriverBindRequest, *gocfbroker.BindingResponse) error
	Unbind(model.DriverUnbindRequest, *string) error
```

## Building and running

### Requirements:
- Go 1.4

### Building:
```sh
mkdir -p $GOPATH/src/github.com/hpcloud
cd $GOPATH/src/github.com/hpcloud
git clone git@github.com:hpcloud/cf-usb.git

make tools
godep restore
make
```
### Running:
To run one or more drivers, you need to copy the driver executable to *{path_to_usb}/drivers/*
```sh
./usb {fileConfigProvider|redisConfigProvider} --{path|uri} {path_to_jsonfile|uri_to_redis}
```