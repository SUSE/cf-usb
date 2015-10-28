package mongo

import "fmt"

type MongoBindingCredentials struct {
	Host             string `json:"server"`
	Port             string `json:"port"`
	Username         string `json:"user_id"`
	Password         string `json:"password"`
	ConnectionString string `json:"connectionString"`
}

var connectionStringTemplate = "mongodb://%[1]v:%[2]v@%[3]v:%[4]v/%[5]v;"

func generateConnectionString(server string, port string, databaseName string, userName string, passWord string) string {
	return fmt.Sprintf(connectionStringTemplate, userName, passWord, server, port, databaseName)
}
