package postgresprovisioner

import (
	"bytes"
	"database/sql"
	//"fmt"
	"log"
	"text/template"

	_ "github.com/lib/pq"
)

var createDatabaseQuery = "CREATE DATABASE {{.Database}} ENCODING 'UTF8'"
var revokeOnDatabaseQuery = "REVOKE all on database {{.Database}} from public"
var createRoleQuery = "CREATE ROLE {{.User}}"
var alterDatabaseOwnerQuery = "ALTER DATABASE {{.Database}} OWNER TO {{.User}}"
var selectCurrentUserQuery = "SELECT current_user"
var terminateDatabaseConnQuery = "SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = {{.Database}} AND pid <> pg_backend_pid()"
var deleteDatabaseQuery = "DROP DATABASE IF EXISTS {{.Database}}"
var deleteRoleQuery = "DROP ROLE IF EXISTS {{.User}}"
var addRoleToDatabaseQuery = "ALTER ROLE {{.User}} LOGIN PASSWORD {{.Password}}"
var alterRoleFromDatabaseQuery = "ALTER ROLE {{.User}} NOLOGIN"

type PostgresProvisioner struct {
	pgClient          *sql.DB
	defaultConnParams map[string]string
}

func NewPostgresProvisioner(defaultConnParams map[string]string) PostgresProvisionerInterface {
	return &PostgresProvisioner{
		pgClient:          nil,
		defaultConnParams: defaultConnParams,
	}
}

func (provisioner *PostgresProvisioner) Init() error {
	var err error = nil
	connString := buildConnectionString(provisioner.defaultConnParams)
	provisioner.pgClient, err = sql.Open("postgres", connString)

	if err != nil {
		return err
	}

	err = provisioner.pgClient.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (provisioner *PostgresProvisioner) Close() error {
	err := provisioner.pgClient.Close()
	return err
}

func (provisioner *PostgresProvisioner) CreateDatabase(dbname string) error {
	err := provisioner.executeQueryNoTx([]string{createDatabaseQuery, revokeOnDatabaseQuery}, map[string]string{"Database": dbname})

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (provisioner *PostgresProvisioner) DeleteDatabase(dbname string) error {

	return nil
}

func (provisioner *PostgresProvisioner) CreateUser(dbname string, username string, password string) error {

	return nil
}

func (provisioner *PostgresProvisioner) DeleteUser(dbname string, username string) error {

	return nil
}

func buildConnectionString(connectionParams map[string]string) string {
	var res string = ""
	for k, v := range connectionParams {
		res += k + "=" + v + ";"
	}
	return res
}

func parametrizeQuery(query string, params map[string]string) string {
	qt := template.Must(template.New("query").Parse(query))
	output := bytes.Buffer{}
	qt.Execute(&output, params)
	return output.String()
}

func (provisioner *PostgresProvisioner) executeQueryNoTx(queries []string, params map[string]string) error {
	for _, query := range queries {
		pQuery := parametrizeQuery(query, params)

		_, err := provisioner.pgClient.Exec(pQuery)
		if err != nil {
			return err
		}
	}

	return nil
}

func (provisioner *PostgresProvisioner) executeQueryTx(queries []string, params map[string]string) error {
	tx, err := provisioner.pgClient.Begin()
	if err != nil {
		return err
	}

	for _, query := range queries {
		pQuery := parametrizeQuery(query, params)

		_, err = tx.Exec(pQuery)
		if err != nil {
			tx.Rollback()

			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (provisioner *PostgresProvisioner) executeQueryRowTx(query string, params map[string]string) error {
	tx, err := provisioner.pgClient.Begin()
	if err != nil {
		return err
	}

	pQuery := parametrizeQuery(query, params)

	var res string

	err = tx.QueryRow(pQuery).Scan(&res)
	if err != nil {
		tx.Rollback()

		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
