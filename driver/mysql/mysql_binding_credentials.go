package driver

import "fmt"

type MysqlBindingCredentials struct {
	Hostname         string `json:"hostname"`
	Host             string `json:"host"`
	Port             string `json:"port"`
	Username         string `json:"username"`
	User             string `json:"user"`
	Password         string `json:"password"`
	Database         string `json:"database"`
	ConnectionString string `json:"connectionString"`
	Uri              string `json:"uri"`
	JdbcUrl          string `json:"jdbcUrl"`
	Name             string `json:"name"`
}

var ConnectionStringTemplate = "Server=%[1]v;Port=%[2]v;Database=%[3]v;Uid=%[4]v;Pwd=%[5]v;"
var UriTemplate = "mysql://%[4]v:%[5]v@%[1]v:%[2]v/%[3]v"
var JdbcUrlTemplate = "jdbc:mysql://%[4]v:%[5]v@%[1]v:%[2]v/%[3]v"

func generateConnections(input string, server string, port string, databaseName string, userName string, passWord string) string {
	return fmt.Sprintf(input, server, port, databaseName, userName, passWord)
}
