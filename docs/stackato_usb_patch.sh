#!/bin/bash

baseurl="https://region-b.geo-1.objects.hpcloudsvc.com/v1/54026737306152/cf-usb-artifacts/verify-141-2016-01-20_14-16-46/linux-amd64.zip"

function random_uuid()
{
# generate a random uuid
echo -n "`cat /dev/urandom | tr -dc 'a-f0-9' | fold -w 8 | head -n 1`-`cat /dev/urandom | tr -dc 'a-f0-9' | fold -w 4 | head -n 1`-`cat /dev/urandom | tr -dc 'a-f0-9' | fold -w 4 | head -n 1`-`cat /dev/urandom | tr -dc 'a-f0-9' | fold -w 4 | head -n 1`-`cat /dev/urandom | tr -dc 'a-f0-9' | fold -w 12 | head -n 1`" 
}

uaa_secret=`random_uuid`
broker_secret=`random_uuid`

function download_bits()
{
# download usb and the drivers
mkdir /s/go/bin/drivers

wget $baseurl

unzip linux-amd64.zip
mv linux-amd64/usb /s/go/bin/
mv linux-amd64/* /s/go/bin/drivers/
rm -rf linux-amd64*
}

function configure_redis()
{
cf_usb_secret=`kato config get aok oauth/clients/cc_usb_management/secret`
aok_verification_key=`kato config get cloud_controller_ng uaa/symmetric_secret`
cc_external_domain=`kato config get cloud_controller_ng external_domain`

echo "Adding management_api to redis"

kato config set usb management_api/listen ":54053"
kato config set usb management_api/uaa_secret "${uaa_secret}"
kato config set usb management_api/uaa_client "cc_usb_management"
kato config set usb management_api/authentication/uaa/adminscope "cloud_controller.admin"
kato config set usb management_api/authentication/uaa/symmetric_verification_key "$aok_verification_key"
kato config set usb management_api/cloud_controller/api "https://$cc_external_domain"
kato config set usb management_api/cloud_controller/skip_tsl_validation "true"
kato config set usb management_api/dev_mode "true"

system_domain=`kato config get cloud_controller_ng system_domain`

echo "Adding broker_api to redis"

kato config set usb broker_api/external_url "http://broker.$system_domain"
kato config set usb broker_api/listen ":54054"
kato config set usb broker_api/credentials/username "usb-broker-admin"
kato config set usb broker_api/credentials/password "${broker_secret}"

system_domain=`kato config get cloud_controller_ng system_domain`
nats_servers=`kato config get -j cloud_controller_ng message_bus_servers|tr -d \\\n|tr -d " "`

echo "Adding routes_register to redis"

kato config set usb routes_register/nats_members "$nats_servers"
kato config set usb routes_register/broker_api_host "broker.$system_domain"

echo "Adding drivers to redis"

for driver in /s/go/bin/drivers/*;
do
        guid=`random_uuid`
        kato config set usb drivers/${guid}/driver_type "`basename $driver`"
        kato config set usb drivers/${guid}/driver_name "`basename $driver`"
done

echo "Adding api_version to redis"

kato config set usb api_version \"2.6\"
}

function configure_role()
{
echo "Configuring role"

cat <<EOF >/s/etc/kato/processes/usb.yml
---
name: usb
roles:
  - usb
EOF

cat <<EOF >>/s/etc/kato/process_order.yml
  -
    name: usb
EOF

cat <<EOF >>/s/etc/kato/role_order.yml
usb:
    min_per_cluster: 0
    exclude_from_add_all: true
EOF

kato config set aok oauth/clients/cc_usb_management/secret "${uaa_secret}"
kato config set aok oauth/clients/cc_usb_management/authorities "cloud_controller.admin"
kato config set aok oauth/clients/cc_usb_management/authorized_grant_types "client_credentials"

cat <<EOF >/s/etc/supervisord.conf.d/usb
[program:usb]
command=/s/go/bin/usb_ctl
priority=5
redirect_stderr=true
stdout_logfile=/s/logs/usb.log
stdout_logfile_maxbytes=1MB
stdout_logfile_backups=3
autostart=false
exitcodes=0
EOF

cat <<"EOF" >/s/go/bin/usb_ctl
#!/bin/bash

function terminate()
{
kill -9 $usbpid
}

trap terminate SIGTERM

redis_uri=`cat /s/etc/kato/redis_uri |cut -f 3 -d \/`
USB_DRIVER_PATH=/s/go/bin/drivers /s/go/bin/usb redisConfigProvider -a ${redis_uri} >>/s/logs/usb.log 2>&1 &
usbpid=$!

wait $usbpid
EOF

chmod +x /s/go/bin/usb_ctl
}

download_bits
configure_role
configure_redis
supervisorctl reread
supervisorctl update
