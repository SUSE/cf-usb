package config

type MssqlDriverConfig struct {
	User               string `json:"userid"`
	Pass               string `json:"password"`
	Host               string `json:"server"`
	Port               int    `json:"port"`
	DbIdentifierPrefix string `json:"db_identifier_prefix"`
}
