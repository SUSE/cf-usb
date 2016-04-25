package config

type MysqlDriverConfig struct {
	User string `json:"userid"`
	Pass string `json:"password"`
	Host string `json:"server"`
	Port string `json:"port"`
}
