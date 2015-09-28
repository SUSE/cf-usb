package mysqlprovisioner

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type MysqlProvisioner struct {
	User string
	Pass string
	Host string
}

func New(username string, password string, host string) MysqlProvisionerInterface {
	return &MysqlProvisioner{User: username, Pass: password, Host: host}
}

func (e *MysqlProvisioner) CreateDatabase(databaseName string) error {

	con, err := e.openSqlConnection()

	defer con.Close()

	if err != nil {
		return err
	}

	err = e.executeTransaction([]string{fmt.Sprintf("CREATE DATABASE %s", databaseName)}, con)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (e *MysqlProvisioner) DeleteDatabase(databaseName string) error {

	con, err := e.openSqlConnection()

	defer con.Close()

	if err != nil {
		return err
	}

	err = e.executeTransaction([]string{fmt.Sprintf("DROP DATABASE %s", databaseName)}, con)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (e *MysqlProvisioner) Query(query string) (*sql.Rows, error) {

	con, err := e.openSqlConnection()

	defer con.Close()

	if err != nil {
		return nil, err
	}

	result, err := con.Query(query)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e *MysqlProvisioner) CreateUser(databaseName string, username string, password string) error {

	con, err := e.openSqlConnection()

	defer con.Close()

	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Connection open - executing transaction")
	err = e.executeTransaction([]string{
		fmt.Sprintf("CREATE USER '%s'@'localhost' IDENTIFIED BY '%s';", username, password),
		fmt.Sprintf("CREATE USER '%s'@'%s' IDENTIFIED BY '%s';", username, e.Host, password),
		fmt.Sprintf("GRANT ALL ON %s.* TO '%s'@'localhost'", databaseName, username)},
		con)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (e *MysqlProvisioner) DeleteUser(username string) error {

	con, err := e.openSqlConnection()

	defer con.Close()

	if err != nil {
		return err
	}
	err = e.executeTransaction([]string{fmt.Sprintf("DROP USER '%s'@'%'", username)}, con)

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
	} else {
		return con, nil
	}
}

func (e *MysqlProvisioner) executeTransaction(querys []string, con *sql.DB) error {
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
