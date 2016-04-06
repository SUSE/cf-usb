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
	Name             string `json:"name"`
	User             string `json:"user"`
	Uri              string `json:"uri"`
	JdbcUrl          string `json:"jdbcUrl"`
}

var connectionStringTemplate = "Server=%[1]v;Port=%[2]v;Database=%[3]v;Uid=%[4]v;Pwd=%[5]v;"
var uriTemplate = "postgres://%[4]v:%[5]v@%[1]v:%[2]v/%[3]v"
var jdbcUrilTemplate = "jdbc:postgresql://%[1]v:%[2]v/%[3]v?user=%[4]v&password=%[5]v"

func generateConnectionString(input string, hostname string, port string, databaseName string, username string, password string) string {
	return fmt.Sprintf(input, hostname, port, databaseName, username, password)
}
