package config

type mynewdriverDriverConfig struct {
	User               string `json:"userid"`
	Pass               string `json:"password"`
	Host               string `json:"host"`
	Port               int    `json:"port"`
}
