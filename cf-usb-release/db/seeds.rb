# This file should contain all the record creation needed to seed the database with its default values.
# The data can then be loaded with the rake db:seed (or created alongside the db with db:setup).
#
# Examples:
#
#   cities = City.create([{ name: 'Chicago' }, { name: 'Copenhagen' }])
#   Mayor.create(name: 'Emanuel', city: cities.first)

connection = ActiveRecord::Base.connection()

sql = <<-EOL
INSERT INTO usb.Config VALUES 
('EXTERNAL_URL','http://1.2.3.4:54053','BROKER_API')
,('LISTEN',':54053','BROKER_API')
,('REQUIRE_TLS','false','BROKER_API')
,('SERVER_CERT_FILE','','BROKER_API')
,('SERVER_KEY_FILE','','BROKER_API')
,('DEV_MODE','true','MANAGEMENT_API')
,('LISTEN',':54054','MANAGEMENT_API')
,('AUTHENTICATION','{"uaa":{"adminscope":"usb.management.admin","public_key":"-----BEGIN PUBLIC KEY-----\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAvExvyYfDFU5tMhgAD+Ak\n9USCMhi6b7sv7EMBIKAy7s83jZUEVR2qgXTQz9hmGy2E7RfA2jz3Jj9XqjIZ38wk\nMWeSHJwaWVvHfSTvx6Gmvax4HK9twUrEo3iSZyTQV7PxzjIGkWuNSE0TwKjTvggU\nTxb7fwyQ9j/x/CIbYVUKLh765seJHuyb4BvbisUQ7l8lBKfFFlrhF2AopneF2P7+\nGYnxz80M8oLMWnYu1c0EQsZ59/E6LYqKvbXM7ZnE9dYguymWQbbWZ06NtcvFFpJN\nTVAT+xhz/ma9R1AMB5gOL84rY2PkmmWzWr4TV4Fe1HWNPzWgZMN9+GNt87AF5/Wr\n9rL8TmHnih/KyV4a/TBCk6pCiyhe/RG+eAhMUaeyNg/a2UZ3kX+OZVF75PAAUeV6\nLA0ZoKyYU9dyP3YqYsaLwUIvCABxCcGVfwmiqrzrSApvck9DKy6U84b3er3GuNw0\n9t1ait99K+YkjU/bJUAkPbdwkt2M5WfdXRT6eN1VBSHIUcb3JjFhuCfosK8tzAmr\n3aZzZ1pdPhXPHNmV0fS8w22L5iavHWWTngLmIF+Ld0bDa5ICgNRykLB7Bcp9lxGF\nmMZAGKVKDL2sCBHeunm68krflztNK8wHCsD/AMeucMJKKf3h6CtW2sSDdNrj/9pU\n6wL/D/IIx+Gd72JtMUuNqzcCAwEAAQ==\n-----END PUBLIC KEY-----"}}','MANAGEMENT_API')
,('UAA_CLIENT','usb-user','MANAGEMENT_API')
,('UAA_SECRET','usb-password','MANAGEMENT_API')
;
INSERT INTO usb.Config VALUES 
('BROKER_NAME','usb','MANAGEMENT_API')
,('API','http://api.192.168.77.77.nip.io','CLOUD_CONTROLLER')
,('SKIP_TLS_VALIDATION','true','CLOUD_CONTROLLER')
,('USERNAME','admin','BROKER_CREDENTIALS')
,('PASSWORD','admin','BROKER_CREDENTIALS')
,('API','1.2.3','API_VERSION')
;
EOL

sql.split(';').each do |s|
  connection.execute(s.strip) unless s.strip.empty?
end

