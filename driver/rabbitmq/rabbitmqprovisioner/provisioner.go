package rabbitmqprovisioner

import (
	"fmt"
	"github.com/streadway/amqp"

	"github.com/hpcloud/cf-usb/driver/rabbitmq/config"
	"github.com/pivotal-golang/lager"
)

type RabbitmqProvisioner struct {
	amqpClient        *amqp.Connection
	defaultConnParams config.RabbitmqDriverConfig
	logger            lager.Logger
}

func NewRabbitmqProvisioner(logger lager.Logger) RabbitmqProvisionerInterface {
	return &RabbitmqProvisioner{logger: logger}
}

func (provisioner *RabbitmqProvisioner) Connect(conf config.RabbitmqDriverConfig) error {
	//var err error = nil
	//connString := buildConnectionString(conf)
	//provisioner.pgClient, err = sql.Open("postgres", connString)

	//if err != nil {
	//	return err
	//}

	//err = provisioner.pgClient.Ping()
	//if err != nil {
	//	return err
	//}

	return nil
}
func (provisioner *RabbitmqProvisioner) Close() error {
	//err := provisioner.pgClient.Close()
	return nil
}

func (provisioner *RabbitmqProvisioner) CreateDatabase(dbname string) error {
	//// for pg driver, create database can not be executed in transaction
	//err := provisioner.executeQueryNoTx([]string{createDatabaseQuery}, map[string]string{"Database": dbname})
	//if err != nil {
	//	return err
	//}

	//err = provisioner.executeQueryTx([]string{revokeOnDatabaseQuery}, map[string]string{"Database": dbname})
	//if err != nil {
	//	return err
	//}

	return nil
}

func (provisioner *RabbitmqProvisioner) DeleteDatabase(dbname string) error {
	//version, err := provisioner.getServerVersion()
	//if err != nil {
	//	return err
	//}

	//var pidColumn string
	//if version > 90200 {
	//	pidColumn = "pid"
	//} else {
	//	pidColumn = "procpid"
	//}

	//err = provisioner.executeQueryTx([]string{terminateDatabaseConnQuery}, map[string]string{
	//	"Database":  dbname,
	//	"PidColumn": pidColumn,
	//})
	//if err != nil {
	//	return err
	//}

	//// for pg driver, drop database can not be executed in transaction
	//err = provisioner.executeQueryNoTx([]string{deleteDatabaseQuery}, map[string]string{"Database": dbname})
	//if err != nil {
	//	return err
	//}

	return nil
}

func (provisioner *RabbitmqProvisioner) DatabaseExists(dbname string) (bool, error) {
	//res, err := provisioner.executeQueryRow(dbCountQuery, map[string]string{"Database": dbname})
	//fmt.Println("count res: ", res)
	//if err != nil {
	//	return false, err
	//}

	//if res.(int64) == 1 {
	//	return true, nil
	//}

	return false, nil
}

func (provisioner *RabbitmqProvisioner) CreateUser(dbname string, username string, password string) error {
	//err := provisioner.executeQueryTx([]string{createRoleQuery, grantAllPrivToRoleQuery}, map[string]string{"User": username, "Password": password, "Database": dbname})
	//if err != nil {
	//	return err
	//}

	return nil
}

func (provisioner *RabbitmqProvisioner) DeleteUser(dbname string, username string) error {
	//err := provisioner.executeQueryTx([]string{revokeAllPrivFromRoleQuery, deleteRoleQuery}, map[string]string{"User": username, "Database": dbname})
	//if err != nil {
	//	return err
	//}

	return nil
}

func (provisioner *RabbitmqProvisioner) UserExists(username string) (bool, error) {
	//res, err := provisioner.executeQueryRow(userCountQuery, map[string]string{"User": username})
	//if err != nil {
	//	return false, err
	//}

	//if res.(int64) == 1 {
	//	return true, nil
	//}

	return false, nil
}

func buildConnectionString(connectionParams config.RabbitmqDriverConfig) string {
	var res string = fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=%v",
		connectionParams.User, connectionParams.Password, connectionParams.Vhost)
	return res
}

//func parametrizeQuery(query string, params map[string]string) (string, error) {
//	queryTemplate := template.Must(template.New("query").Parse(query))
//	output := bytes.Buffer{}
//	queryTemplate.Execute(&output, params)

//	queryString := output.String()

//	if strings.Contains(queryString, "<no value>") {
//		return queryString, errors.New("Invalid parameter passed to query")
//	}

//	return queryString, nil
//}

//func (provisioner *PostgresProvisioner) executeQueryNoTx(queries []string, params map[string]string) error {
//	for _, query := range queries {
//		pQuery, err := parametrizeQuery(query, params)

//		provisioner.logger.Debug("postgres-exec", lager.Data{"query": pQuery})
//		if err != nil {
//			provisioner.logger.Error("postgres-exec", err, lager.Data{"query": pQuery})
//			return err
//		}

//		_, err = provisioner.pgClient.Exec(pQuery)
//		if err != nil {
//			provisioner.logger.Error("postgres-exec", err, lager.Data{"query": pQuery})
//			return err
//		}
//	}

//	return nil
//}

//func (provisioner *PostgresProvisioner) executeQueryTx(queries []string, params map[string]string) error {
//	tx, err := provisioner.pgClient.Begin()
//	if err != nil {
//		return err
//	}

//	for _, query := range queries {
//		pQuery, err := parametrizeQuery(query, params)
//		provisioner.logger.Debug("postgres-exec", lager.Data{"query": pQuery})
//		if err != nil {
//			provisioner.logger.Error("postgres-exec", err, lager.Data{"query": pQuery})
//			return err
//		}

//		_, err = tx.Exec(pQuery)
//		if err != nil {
//			tx.Rollback()
//			provisioner.logger.Error("postgres-exec", err, lager.Data{"query": pQuery})
//			return err
//		}
//	}

//	err = tx.Commit()
//	if err != nil {
//		return err
//	}

//	return nil
//}

//func (provisioner *PostgresProvisioner) executeQueryRow(query string, params map[string]string) (interface{}, error) {
//	pQuery, err := parametrizeQuery(query, params)
//	provisioner.logger.Debug("postgres-exec", lager.Data{"query": pQuery})
//	if err != nil {
//		provisioner.logger.Error("postgres-exec", err, lager.Data{"query": pQuery})
//		return nil, err
//	}

//	var res interface{}
//	err = provisioner.pgClient.QueryRow(pQuery).Scan(&res)
//	if err != nil && err == sql.ErrNoRows {
//		provisioner.logger.Error("postgres-exec", err, lager.Data{"query": pQuery})
//		return nil, err
//	}

//	return res, nil
//}

func (provisioner *RabbitmqProvisioner) getServerVersion() (int, error) {
	//res, err := provisioner.executeQueryRow("SHOW server_version_num", map[string]string{})
	//if err != nil {
	//	return 0, err
	//}

	//i := res.([]uint8)
	//b := make([]byte, len(i))
	//for i, v := range i {
	//	if v < 0 {
	//		b[i] = byte(256 + int(v))
	//	} else {
	//		b[i] = byte(v)
	//	}
	//}

	//return strconv.Atoi(string(b))

	return 0, nil
}
