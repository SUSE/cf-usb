# Universal Service Broker

## Summary
The cf-usb project implements and exposes the Cloud Foundry [Service Broker API](http://docs.cloudfoundry.org/services/api.html). 

It uses plugins (drivers) to connect to different services.

### Related Projects

* [cf-usb-plugin](https://github.com/SUSE/cf-usb-plugin)
* [cf-usb-sidecar](https://github.com/SUSE/cf-usb-sidecar)

## Configuration and Usage

The Universal Service Broker has three main components:

* Universal Service Broker
* Sidecars
* cf CLI plugin

The USB runs as a component within Cloud Foundry. Its [BOSH release](https://github.com/SUSE/cf-usb/tree/develop/cf-usb-release)
is inside the USB main repo. This is included as a submodule in the SUSE CF project, and runs by default.

These usage instructions assume:
* A working SUSE CAP instance is available
* The SUSE CAP instance is able to connect to a Docker registry
* The user is able to access the Kubernetes cluster with the `kubectl` tool
* The user has admin access to the SUSE CAP cluster

### Sidecar Setup

The USB itself is just a broker, and doesn't run any actual services. These are
provided by the [sidecars](https://github.com/SUSE/cf-usb-sidecar/tree/develop/csm-extensions/services), and run outside of the CF cluster.

To build the sidecar, check out the sidecar project:

```
mkdir -p $GOPATH/src/github.com/SUSE
git clone https://github.com/SUSE/cf-usb-sidecar $GOPATH/src/github.com/SUSE/cf-usb-sidecar
cd $GOPATH/src/github.com/SUSE/cf-usb-sidecar
```

Configure your Docker repository information:

```
export DOCKER_REPOSITORY=docker.io
export DOCKER_ORGANIZATION=splatform

docker login # if authorization is required
```

Then build the top level dependencies and the sidecar:

```
make tools
make build-image
cd csm-extensions/services/dev-mysql
make build-image publish-image helm
```

The generated Helm chart will be available in the `chart/` directory.

Install the Helm chart:

```
# You will need to know the namespaces and domain for your cluster:
UAA_NAMESPACE=uaa
CF_NAMESPACE=cf
CF_DOMAIN=cf-dev.io

SIDECAR_NAMESPACE=mysql

UAA_CA_CERT="$(kubectl get secret secret --namespace ${UAA_NAMESPACE} -o jsonpath="{.data['internal-ca-cert']}" | base64 --decode -)"
CF_CA_CERT="$(kubectl get secret secret --namespace ${CF_NAMESPACE} -o jsonpath="{.data['internal-ca-cert']}" | base64 --decode -)"
CF_PASSWORD="$(kubectl get secret secret --namespace ${CF_NAMESPACE} -o jsonpath="{.data['cluster-admin-password']}" | base64 --decode -)"

helm install ./chart --name mysql-instance --namespace ${SIDECAR_NAMESPACE} \
	--set "env.UAA_CA_CERT=${UAA_CA_CERT}" \
	--set "env.CF_CA_CERT=${CF_CA_CERT}" \
	--set "env.SERVICE_LOCATION=http://cf-usb-sidecar-mysql.${SIDECAR_NAMESPACE}.svc.cluster.local:8081" \
	--set "env.SERVICE_MYSQL_HOST=AUTO" \
	--set "env.CF_ADMIN_USER=admin" \
	--set "env.CF_ADMIN_PASSWORD=${CF_PASSWORD}" \
	--set "env.CF_DOMAIN=${CF_DOMAIN}"
```

Eventually you should see two pods start in the `SIDECAR_NAMESPACE`:

```
$ kubectl get pods --namespace ${SIDECAR_NAMESPACE}
NAME                                      READY     STATUS    RESTARTS   AGE
cf-usb-sidecar-mysql-b44d4d66f-d27qb      1/1       Running   0          2m
mysql-0                                   1/1       Running   0          2m
```

There will also be a 'setup' pod that starts an errand, but it will finish and exit.

Once the pods are ready, it should be available in the marketplace:

```
$ cf marketplace
Getting services from marketplace in org org / space space as admin...
OK

service    plans     description
postgres   default   Default service
mysql      default   Default service

TIP:  Use 'cf marketplace -s SERVICE' to view descriptions of individual plans of a given service.
```

At this point, services can be made available to apps. In this case we're going to use the [django-cms](https://github.com/scf-samples/django-cms) app.

```
git clone https://github.com/scf-samples/django-cms -b scf
cd django-cms
cf create-service postgres default django-cms-db
cf push --no-start django-cms
cf set-env django-cms DISABLE_COLLECTSTATIC 1
cf set-env django-cms DJANGO_SETTINGS_MODULE settings
cf start django-cms
```

## Unmanaged USB

Unmanaged USB provides provides a limited number of features and implements the basic functionality of a Cloud Foundry Service Broker.

Constraints:
- it does not provide functionality for upgrading drivers, services, plans, etc.
- does not automatically register to the Cloud Controller.
- only `fileConfigProvider` can be used as a configuration provider.
- it does not expose a management API.
- if a `driver` fails to start, the connection between the `driver` and the server cannot be establised or if the configuration/dials schema can not be validated, the USB exits with an exitcode != 0.

### Deployment strategies

#### 1. Cloud Foundry application

The unmanaged USB can be deployed as an app to a Cloud Foundry deployment. All the services managed must be external services.

#### 2. fissile

TODO.

## Managed USB

The managed USB provides a management API for configuration and update.

Constraints:
- it can not use `fileConfigProvider` as a configuration provider

### Management API

#### Authorization

##### 1. UAA
USB uses UAA as an authorization provider. It requires the `cc_usb_management` OAuth client to be configured with the following properties:
- secret: {clientsecret}
- scope: cloud_controller.write,openid,cloud_controller.read,cloud_controller_service_permissions.read, cloud_controller_service_permissions.write
- authorities: usb.management.admin
- authorized-grant-types: client_credentials 

##### 2. Basic auth
Basic auth can be used when making calls to the USB management API.

#### API Definition

The usb management API is described [here](https://github.com/SUSE/cf-usb/blob/b84f846eedc13c2cf9215c53f323b01c545aab42/docs/mgmt.html)

### fissile

TODO:

## Configuration

**cf-usb** works with multiple configuration providers:

### File configuration provider
**cf-usb** can take it's configuration from a .json file. The json file format is similar to a standard broker configuration file. 
To start the broker with a file configuration provider you must provide the following cli options:

```sh
./usb fileConfigProvider --path {path_to_jsonfile}
```

`dials`

`driver_configs`

### MySQL configuration provider
The state and the configuration of USB can be stored in a MySQL database.
To start the broker using a MySQL provider you must provide the following cli options:

| Option              | Description |
| ------------------- | --- |
| `--address`,  `-a`  | server address and port (mandatory) |
| `--database`, `-db` | database name |
| `--username`, `-u`  | username |
| `--password`, `-p`  | password |

## Drivers

### Folder structure
```
|-- cmd
|	|-- main.go
|-- drivers
|   |-- {driver_type}
|       |-- data
|           |-- schemas.go     | Generated by usb Makefile
|       |-- schemas
|           |-- dails.json     | dails JSON schema definition
|           |-- config.json    | driver_config JSON schema definition

```

### Principals

- Each driver call must be represented by an atomic operation.
- USB must be able to validate all incoming CC requests using the driver interface.

### Driver Interface

### Ping
Checks if the driver can reach the server.

**Request** 

| Type | Description |
| :---: | :---: |
| *json.RawMessage | Driver configuration object |


**Response**

| Type | Description |
| :---: | :---: |
| bool | Returns *true* if the driver can reach the server, *false* otherwise |
| error | Service connection error |

### GetDialsSchema

**Request** 

| Type | Description |
| :---: | :---: |
| string | empty string |

**Response**

| Type | Description |
| :---: | :---: |
| string | the json schema for the *dials* supported by the driver |
| error | GetDialsSchema error |

Example schema for MSSQL *dials*
```sh
{
"type": "object",
"properties": {
"max_dbsize_mb": {
"type": "integer",
"minimum": 0
}
},
"required": ["max_dbsize_mb"]
}
```

### GetConfigSchema

Gets the JSON schema for the *driver_config*

**Request** 

| Type | Description |
| :---: | :---: |
| string | empty string |

**Response**

| Type | Description |
| :---: | :---: |
| string | the json schema for the *driver_config* of the driver |
| error | GetDialsSchema error |

Example schema for MSSQL *driver_config*
```sh
{
"type": "object",
"properties": {
"brokerGoSqlDriver": {
"type": "string"
},
"brokerMssqlConnection": {
"type": "object",
"properties": {
"server": {
"type": "string"
},
"port": {
"type": "string"
},
"database": {
"type": "string"
},
"user id": {
"type": "string"
},
"password": {
"type": "string"
}
}
}
},
"required": [
"brokerGoSqlDriver",
"brokerMssqlConnection"
]
}
```

### ProvisionInstance
Creates a service instance

**Request**

| Type | Description |
| :---: | :---: |
| ProvisionInstanceRequest | Object containing the *instanceID*, *config* object and  *dails* object |

Dials are restrictions that can be applied to instances.
Example for MSSQL:
```sh
{
"max_dbsize_mb": 1500
}
```

**Response**

| Type | Description |
| :---: | :---: |
| Instance | Instance type object containing *instanceID*, *status* and *description* |
| error | Provision Instance error |

### GetInstance

Retrieves an existing instance.

**Request**

| Type | Description |
| :---: | :---: |
| GetInstanceRequest | Object containing the *instanceID* and a *config* object |
**Response**

| Type | Description |
| :---: | :---: |
| Instance | Instance type object containing *instanceID*, *status* and *description*  |
| error | Instance Exists error |

### GenerateCredentials

Generates for the for the specified instance

**Request**

| Type | Description |
| :---: | :---: |
| GenerateCredentialsRequest | Object containing the *instanceID*, the *credentialsID* and the *config* object |

**Response**

| Type | Description |
| :---: | :---: |
| interface{} | Connection information for the specified service |
| error | GenerateCredentials error |

Example for MSSQL:

```sh
type MssqlCredentials struct {
Hostname         string `json:"hostname"`
Port             int    `json:"port"`
Name             string `json:"name"`
Username         string `json:"username"`
Password         string `json:"password"`
ConnectionString string `json:"connectionString"`
}
```

### GetCredentials

Retrieves existing credentials.

**Request**

| Type | Description |
| :---: | :---: |
| GetCredentialsRequest | Object containing the *instanceID*, the *credentialsID* and the *config* object |
**Response**

| Type | Description |
| :---: | :---: |
| Credentials | Credentials type object containing *credentialsID*, *status* and *description* |
| error | CredentialsExists error |

### RevokeCedentials

Revoke the credentials of the specified *credentialsID*

**Request**

| Type | Description |
| :---: | :---: |
| RevokeCredentialsRequest | Object containing the *instanceID*, the *credentialsID* and the *config* object |
**Response**

| Type | Description |
| :---: | :---: |
| Credentials | Credentials type object containing *credentialsID*, *status* and *description* |
| error | RevokeCedentials error |


### DeprovisionInstance

**Request**

Deprovisions the instance having the specified *instanceID*

| Type | Description |
| :---: | :---: |
| DeprovisionInstanceRequest | Object containing the *instanceID* and the *config* object |

**Response**

| Type | Description |
| :---: | :---: |
| Instance | Instance type object containing *instanceID*, *status* and *description*  |
| error | DeprovisionInstance error |



## Building and running

### Requirements:
- Go 1.6

### Building:
```sh
mkdir -p $GOPATH/src/github.com/SUSE
cd $GOPATH/src/github.com/SUSE
git clone git@github.com:SUSE/cf-usb.git

make tools
godep restore
make
```
### Running:
To run one or more drivers, you need to copy the driver executable to *{path_to_usb}/drivers/* or set the USB_DRIVER_PATH environment variable with the path to the driver executable 
```sh
./usb {ConfigProvider} --{options}
```
