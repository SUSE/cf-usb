package config

type MongoDriverConfig struct {
	User string `json:"userid"`
	Pass string `json:"password"`
	Host string `json:"server"`
	Port string `json:"port"`
}
