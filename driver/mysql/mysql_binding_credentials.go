package driver

import "fmt"

type MysqlBindingCredentials struct {
	Host             string `json:"server"`
	Port             string `json:"port"`
	Username         string `json:"user_id"`
	Password         string `json:"password"`
	Database         string `json:"database"`
	ConnectionString string `json:"connectionString"`
}

var connectionStringTemplate = "Server=%[1]v;Port=%[2]v;Database=%[3]v;Uid=%[4]v;Pwd=%[5]v;"

func generateConnectionString(server string, port string, databaseName string, userName string, passWord string) string {
	return fmt.Sprintf(connectionStringTemplate, server, port, databaseName, userName, passWord)
}
