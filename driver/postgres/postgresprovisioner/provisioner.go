package postgresprovisioner

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"text/template"

	_ "github.com/lib/pq"
	"github.com/pivotal-golang/lager"
)

var createDatabaseQuery = "CREATE DATABASE {{.Database}} ENCODING 'UTF8'"
var revokeOnDatabaseQuery = "REVOKE all on database {{.Database}} from public"
var createRoleQuery = "CREATE ROLE {{.User}} LOGIN PASSWORD '{{.Password}}'"
var grantAllPrivToRoleQuery = "GRANT ALL PRIVILEGES ON DATABASE {{.Database}} TO {{.User}}"
var revokeAllPrivFromRoleQuery = "REVOKE ALL PRIVILEGES ON DATABASE {{.Database}} FROM {{.User}}"
var deleteRoleQuery = "DROP ROLE {{.User}}"
var terminateDatabaseConnQuery = "SELECT pg_terminate_backend(pg_stat_activity.procpid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '{{.Database}}' AND procpid <> pg_backend_pid()"
var deleteDatabaseQuery = "DROP DATABASE {{.Database}}"

type PostgresServiceProperties struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Dbname   string `json:"dbname"`
	Sslmode  string `json:"sslmode"`
}

type PostgresProvisioner struct {
	pgClient          *sql.DB
	defaultConnParams PostgresServiceProperties
	logger            lager.Logger
}

func NewPostgresProvisioner(defaultConnParams PostgresServiceProperties, logger lager.Logger) PostgresProvisionerInterface {
	return &PostgresProvisioner{
		pgClient:          nil,
		defaultConnParams: defaultConnParams,
		logger:            logger,
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
		return err
	}

	return nil
}

func (provisioner *PostgresProvisioner) DeleteDatabase(dbname string) error {
	err := provisioner.executeQueryTx([]string{terminateDatabaseConnQuery}, map[string]string{"Database": dbname})
	if err != nil {
		return err
	}

	err = provisioner.executeQueryNoTx([]string{deleteDatabaseQuery}, map[string]string{"Database": dbname})
	if err != nil {
		return err
	}

	return nil
}

func (provisioner *PostgresProvisioner) CreateUser(dbname string, username string, password string) error {
	err := provisioner.executeQueryTx([]string{createRoleQuery, grantAllPrivToRoleQuery}, map[string]string{"User": username, "Password": password, "Database": dbname})
	if err != nil {
		return err
	}

	return nil
}

func (provisioner *PostgresProvisioner) DeleteUser(dbname string, username string) error {
	err := provisioner.executeQueryTx([]string{revokeAllPrivFromRoleQuery, deleteRoleQuery}, map[string]string{"User": username, "Database": dbname})
	if err != nil {
		return err
	}

	return nil
}

func buildConnectionString(connectionParams PostgresServiceProperties) string {
	var res string = fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=%v",
		connectionParams.User, connectionParams.Password, connectionParams.Host, connectionParams.Port, connectionParams.Sslmode)

	return res
}

func parametrizeQuery(query string, params map[string]string) (string, error) {
	queryTemplate := template.Must(template.New("query").Parse(query))
	output := bytes.Buffer{}
	queryTemplate.Execute(&output, params)

	queryString := output.String()

	if strings.Contains(queryString, "<no value>") {
		return queryString, errors.New("Invalid parameter passed to query")
	}

	return queryString, nil
}

func (provisioner *PostgresProvisioner) executeQueryNoTx(queries []string, params map[string]string) error {
	for _, query := range queries {
		pQuery, err := parametrizeQuery(query, params)

		provisioner.logger.Debug("postgres-exec", lager.Data{"query": pQuery})
		if err != nil {
			provisioner.logger.Error("postgres-exec", err, lager.Data{"query": pQuery})
			return err
		}

		_, err = provisioner.pgClient.Exec(pQuery)
		if err != nil {
			provisioner.logger.Error("postgres-exec", err, lager.Data{"query": pQuery})
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
		pQuery, err := parametrizeQuery(query, params)
		provisioner.logger.Debug("postgres-exec", lager.Data{"query": pQuery})
		if err != nil {
			provisioner.logger.Error("postgres-exec", err, lager.Data{"query": pQuery})
			return err
		}

		_, err = tx.Exec(pQuery)
		if err != nil {
			tx.Rollback()
			provisioner.logger.Error("postgres-exec", err, lager.Data{"query": pQuery})
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
