In the following, we're going to suppose that the public ip is `54.154.46.52` and the domain is `helion-cf.io`

#### 1. Add usb roule to AOK
In the redis database, edit the `aok` key and add the following (json format) :
```
  "cc_usb_management": {
                "secret": "cc-usb-management-secret",
                "authorities": "cloud_controller.admin",
                "authorized_grant_types": "client_credentials",
                "scope": "usb.management.admin"
            },
```

```
kato config set aok oauth/clients/cc_usb_management/secret cc-usb-management-secret
kato config set aok oauth/clients/cc_usb_management/authorities cloud_controller.admin,usb.management.admin
kato config set aok oauth/clients/cc_usb_management/authorized_grant_types client_credentials
kato config set aok oauth/clients/cc_usb_management/scope usb.management.admin
```

#### 2. Add *usb.management.admin* oauth-> users -> default_authorities
```
kato config push aok oauth/users/default_authorities usb.management.admin
```
#### 3. Get cloud_controller.api from *cloud_controller_ng/external_domain*
(use this to validate the initial assumptions)
```
kato config get cloud_controller_ng external_domain
```
#### 4. Get public key
You will use this at step #14.
```
kato config get cloud_controller_ng uaa/symmetric_secret
```
#### 5. Create share between USB drivers
You can do this using nfs mount. Use the env var USB_DRIVER_PATH to specify the shared directory.
Or, alternatively you can use sshfs (assuming you have the drivers on 54.154.46.53)
```
mkdir -p /s/go/bin/drivers
sshfs -p 22 -o idmap=user -o reconnect -o ServerAliveInterval=15 stackato@54.154.46.53:/s/go/bin/drivers /s/go/bin/drivers
```
#### 6. Advertise broker.DOMAIN and usb.DOMAIN if a node is attached
(replace DOMAIN with a proper domain)
```
gem install nats
nats-pub 'router.register' '{"host":"127.0.0.1","port":54053,"uris":["usb.54.154.36.52.helion-cf.io"]}'
nats-pub 'router.register' '{"host":"127.0.0.1","port":54054,"uris":["broker.54.154.36.52.helion-cf.io"]}'
```
#### 7. Add to /s/etc/kato/role_order.yml
```
usb:
    min_per_cluster: 0
    max_per_cluster: 1
    exclude_from_add_all: true
```
#### 8. Add to /s/etc/kato/process_order.yml
```
  -
    name: usb
```
#### 9. Add supervisord config to: /s/etc/supervisord.conf.d/usb
```
[program:usb]
command=USB_DRIVER_PATH=/s/go/bin/drivers /s/go/bin/usb redisConfigProvider -a 127.0.0.1:7474
priority=5
redirect_stderr=true
stdout_logfile=/s/logs/usb.log
stdout_logfile_maxbytes=1MB
stdout_logfile_backups=3
autostart=false
exitcodes=0
```
#### 10. usb executable should be copied to : /s/go/bin
These  can be found on swift, tenant `hpcs-apaas-tenant1`, zone `us-east`, container `cf-usb-artifacts`. For now you can use the artifacts made by the `verify` job
#### 11. usb drivers should be copied to : /s/go/bin/drivers
These  can be found on swift, tenant `hpcs-apaas-tenant1`, zone `us-east`, container `cf-usb-artifacts`. For now you can use the artifacts made by the `verify` job
#### 12. Create /s/etc/kato/processes/usb.yml
```
---
name: usb
roles:
  - usb
```
#### 13. Add "management_api" and "broker_api" keys to redis (port 7474)
For now, we use "dev_mode=true" to skip token authentication as it is not yet fully implemented for stackato. Your public key will not be taken into consideration if "dev_mode=true".
```
set management_api '{"dev_mode":true,"listen":":54053","uaa_secret":"cc-usb-management-secret","uaa_client":"cc_usb_management","authentication":{"uaa":{"adminscope":"usb.management.admin","public_key":"YOUR PUBLIC KEY GOES HERE"}},"cloud_controller":{"api":"https://api.54.154.46.52.helion-cf.io","skip_tsl_validation":true}}'
```
```
set broker_api '{"external_url":"http://broker.54.154.46.52.helion-cf.io","listen":":54054","credentials":{"username":"username","password":"password"}}'
```
#### 14. Start USB
You can start usb either with supervisord or manually running :
```
USB_DRIVER_PATH=/s/go/bin/drivers /s/go/bin/usb redisConfigProvider -a 127.0.0.1:7474
```