package redis

type RedisBindingCredentials struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}
