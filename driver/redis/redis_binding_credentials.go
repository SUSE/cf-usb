package redis

type RedisBindingCredentials struct {
	Hostname string `json:"hostname"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}
