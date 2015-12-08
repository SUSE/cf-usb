package config

type RabbitmqDriverConfig struct {
	Vhost    string `json:"vhost"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}
