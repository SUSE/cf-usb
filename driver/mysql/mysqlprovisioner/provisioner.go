package mysqlprovisioner

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type MysqlProvisioner struct {
	User       string
	Pass       string
	Host       string
	Connection *sql.DB
}

func New(username string, password string, host string) (MysqlProvisionerInterface, error) {
	var err error
	provisioner := MysqlProvisioner{User: username, Pass: password, Host: host}

	provisioner.Connection, err = provisioner.openSqlConnection()

	if err != nil {
		return nil, err
	}

	return &provisioner, nil
}

func (e *MysqlProvisioner) Close() error {
	err := e.Connection.Close()
	return err
}

func (e *MysqlProvisioner) IsDatabaseCreated(databaseName string) (bool, error) {
	rows, err := e.Query("SHOW DATABASES")
	if err != nil {
		return false, err
	}

	var (
		result    [][]string
		container []string
		pointers  []interface{}
	)

	cols, err := rows.Columns()
	if err != nil {
		return false, err
	}

	length := len(cols)

	for rows.Next() {
		pointers = make([]interface{}, length)
		container = make([]string, length)

		for i := range pointers {
			pointers[i] = &container[i]
		}

		err = rows.Scan(pointers...)
		if err != nil {
			return false, err
		}

		result = append(result, container)
	}
	for _, cont := range result {
		if cont[0] == databaseName {
			return true, nil
		}
	}

	return false, nil
}

func (e *MysqlProvisioner) IsUserCreated(userName string) (bool, error) {
	rows, err := e.Query("SELECT user from mysql.user")
	if err != nil {
		return false, err
	}

	var (
		result    [][]string
		container []string
		pointers  []interface{}
	)

	cols, err := rows.Columns()
	if err != nil {
		return false, err
	}

	length := len(cols)

	for rows.Next() {
		pointers = make([]interface{}, length)
		container = make([]string, length)

		for i := range pointers {
			pointers[i] = &container[i]
		}

		err = rows.Scan(pointers...)
		if err != nil {
			return false, err
		}

		result = append(result, container)
	}
	for _, cont := range result {
		if cont[0] == userName {
			return true, nil
		}
	}

	return false, nil
}

func (e *MysqlProvisioner) CreateDatabase(databaseName string) error {
	err := e.executeTransaction(e.Connection, fmt.Sprintf("CREATE DATABASE %s", databaseName))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (e *MysqlProvisioner) DeleteDatabase(databaseName string) error {

	err := e.executeTransaction(e.Connection, fmt.Sprintf("DROP DATABASE %s", databaseName))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (e *MysqlProvisioner) Query(query string) (*sql.Rows, error) {

	result, err := e.Connection.Query(query)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e *MysqlProvisioner) CreateUser(databaseName string, username string, password string) error {

	log.Println("Connection open - executing transaction")
	err := e.executeTransaction(e.Connection,
		fmt.Sprintf("CREATE USER '%s'@'localhost' IDENTIFIED BY '%s';", username, password),
		fmt.Sprintf("CREATE USER '%s'@'%s' IDENTIFIED BY '%s';", username, e.Host, password),
		fmt.Sprintf("GRANT ALL ON %s.* TO '%s'@'localhost'", databaseName, username))
	log.Println("Transaction done")
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (e *MysqlProvisioner) DeleteUser(username string) error {

	err := e.executeTransaction(e.Connection, fmt.Sprintf("DROP USER '%s'@'localhost'", username),
		fmt.Sprintf("DROP USER '%s'@'%s'", username, e.Host))

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (e *MysqlProvisioner) openSqlConnection() (*sql.DB, error) {
	con, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/mysql", e.User, e.Pass, e.Host))
	if err != nil {
		return nil, err
	}
	return con, nil
}

func (e *MysqlProvisioner) executeTransaction(con *sql.DB, querys ...string) error {
	tx, err := con.Begin()
	if err != nil {
		log.Println(err)
		return err
	} else {
		for _, query := range querys {
			_, err = tx.Exec(query)
			if err != nil {
				log.Println(err)
				tx.Rollback()
				break
			}
		}
		tx.Commit()
	}
	return nil
}
