package postgres

import (
	"fmt"
)

type PostgresBindingCredentials struct {
	Hostname         string `json:"hostname"`
	Host             string `json:"host"`
	Database         string `json:"database"`
	Password         string `json:"password"`
	Port             string `json:"port"`
	Username         string `json:"username"`
	ConnectionString string `json:"connectionString"`
}

var connectionString = "Server=%[1]v;Port=%[2]v;Database=%[3]v;Uid=%[4]v;Pwd=%[5]v;"

func generateConnectionString(hostname string, port string, databaseName string, username string, password string) string {
	return fmt.Sprintf(connectionString, hostname, port, databaseName, username, password)
}
