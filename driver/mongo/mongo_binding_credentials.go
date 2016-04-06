package mongo

import "fmt"

type MongoBindingCredentials struct {
	Hostname string `json:"hostname"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Uri      string `json:"uri"`
	Name     string `json:"name"`
	Db       string `json:"db"`
}

var connectionStringTemplate = "mongodb://%[1]v:%[2]v@%[3]v:%[4]v/%[5]v;"

func generateConnectionString(server string, port string, databaseName string, userName string, passWord string) string {
	return fmt.Sprintf(connectionStringTemplate, userName, passWord, server, port, databaseName)
}
