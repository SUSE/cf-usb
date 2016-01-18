In the following, we're going to suppose that the public ip is `54.154.46.52` and the domain is `helion-cf.io`

#### 1. Add usb role to AOK
In the redis database, edit the `aok` key and add the following (json format) :
```
        "cc_usb_management": {
                "secret": "CHANGE-ME-SECRET",
                "authorities": "cloud_controller.admin",
                "authorized_grant_types": "client_credentials"
        }
```

```
kato config set aok oauth/clients/cc_usb_management/secret "CHANGE-ME-SECRET"
kato config set aok oauth/clients/cc_usb_management/authorities "cloud_controller.admin"
kato config set aok oauth/clients/cc_usb_management/authorized_grant_types "client_credentials"
```
#### 2. Create share between USB drivers
You can do this using nfs mount. Use the env var USB_DRIVER_PATH to specify the shared directory.
Or, alternatively you can use sshfs (assuming you have the drivers on 54.154.46.53)
```
mkdir -p /s/go/bin/drivers
sshfs -p 22 -o idmap=user -o reconnect -o ServerAliveInterval=15 stackato@54.154.46.53:/s/go/bin/drivers /s/go/bin/drivers
```
#### 3. Add to /s/etc/kato/role_order.yml
```
usb:
    min_per_cluster: 0
    exclude_from_add_all: true
```
#### 4. Add to /s/etc/kato/process_order.yml
```
  -
    name: usb
```
#### 5. Add supervisord config to: /s/etc/supervisord.conf.d/usb
```
[program:usb]
command=USB_DRIVER_PATH=/s/go/bin/drivers /s/go/bin/usb redisConfigProvider -a {core_ip}:7474
priority=5
redirect_stderr=true
stdout_logfile=/s/logs/usb.log
stdout_logfile_maxbytes=1MB
stdout_logfile_backups=3
autostart=false
exitcodes=0
```
#### 6. usb executable should be copied to : /s/go/bin
These  can be found on swift, tenant `hpcs-apaas-tenant1`, zone `us-east`, container `cf-usb-artifacts`. For now you can use the artifacts made by the `verify` job
#### 7. usb drivers should be copied to : /s/go/bin/drivers
These  can be found on swift, tenant `hpcs-apaas-tenant1`, zone `us-east`, container `cf-usb-artifacts`. For now you can use the artifacts made by the `verify` job
#### 8. Create /s/etc/kato/processes/usb.yml
```
---
name: usb
roles:
  - usb
```
#### 9. Add "api_version", "management_api", "broker_api", "drivers" and "routes_register" keys to redis (port 7474)
Set "api_version"
```
redis-cli -p 7474 set api_version "2.0"
```

Set  "management_api"
```
export cf_usb_secret=`kato config get aok oauth/clients/cc_usb_management/secret`
export aok_verification_key=`kato config get cloud_controller_ng uaa/symmetric_secret`
export cc_external_domain=`kato config get cloud_controller_ng external_domain`

read -d '' management_api << EOF
{
        "listen": ":54053",
        "uaa_secret": "$cf_usb_secret",
        "uaa_client": "cc_usb_management",
        "authentication": {
                "uaa": {
                        "adminscope": "cloud_controller.admin",
                        "symmetric_verification_key": "$aok_verification_key"
                }
        },
        "cloud_controller": {
                "api": "https://$cc_external_domain",
                "skip_tsl_validation": true
        }
}
EOF

redis-cli -p 7474 set management_api "$management_api"
```


Set "broker_api" (n.b.)
```
export system_domain=`kato config get cloud_controller_ng system_domain`

read -d '' broker_api << EOF
{
    "external_url": "http://broker.$system_domain",
    "listen": ":54054",
    "credentials": {
        "username": "usb-broker-admin",
        "password": "CHANGE-ME-PASSWORD"
    }
}
EOF

redis-cli -p 7474 set broker_api "$broker_api"
```

Set "routes_register"
```
export system_domain=`kato config get cloud_controller_ng system_domain`
export nats_servers=`kato config get -j cloud_controller_ng message_bus_servers`

read -d '' routes_register << EOF
{
    "nats_members": $nats_servers,
    "broker_api_host": "broker.$system_domain",
    "management_api_host": "management.$system_domain"
}
EOF

redis-cli -p 7474 set routes_register "$routes_register"
```

For the drivers key, you need to compose a json string which contains references to all drivers. In the end it should look like this :
```
read -d '' drivers << EOF
{
  "df652dab-478e-3020-30b8-59054fc0bd6e": {
    "driver_type": "dummy"
  },
  "67d8cc5d-5fdf-d199-c184-b01701875773": {
    "driver_type": "mongo"
  },
  "5921f46b-95a9-bb6f-9365-040ca5afd69e": {
    "driver_type": "mssql"
  },
  "7d991a4d-6b1f-de18-e22d-0c92860530a7": {
    "driver_type": "mysql"
  },
  "14a067c9-06af-2c59-5c92-df7d02cf592c": {
    "driver_type": "postgres"
  },
  "d4d74653-6a56-8190-2e17-06bd894319c7": {
    "driver_type": "rabbitmq"
  },
  "63a67bde-152d-2195-3b0d-f62b1db1b2d3": {
    "driver_type": "redis"
  }
}
EOF
redis-cli -p 7474 set drivers "$drivers"

```

You can use the following script to generate the json string :
```
#!/bin/bash 

function random_uuid() 
{ 
   echo -n "`cat /dev/urandom | tr -dc 'a-f0-9' | fold -w 8 | head -n 1`-`cat /dev/urandom | tr -dc 'a-f0-9' | fold -w 4 | head -n 1`-`cat /dev/urandom | tr -dc 'a-f0-9' |
fold -w 4 | head -n 1`-`cat /dev/urandom | tr -dc 'a-f0-9' | fold -w 4 | head -n 1`-`cat /dev/urandom | tr -dc 'a-f0-9' | fold -w 12 | head -n 1`" 
} 


function generate_json() 
{ 
echo -n "{" 

for driver in /s/go/bin/drivers/*; 
       do 
       echo -n "\"`random_uuid`\": { \"driver_type\": \"`basename $driver`\" }," 
       done 
} 

final=`generate_json|rev|cut -b 2-|rev` 
echo "${final}}"
```

#### 10. Restart all
```
kato restart
```
#### 11. Restart supervisord
```
stop-supervisord
start-supervisord
```
