package mongoprovisioner

import (
	"fmt"
	"github.com/pivotal-golang/lager"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoProvisioner struct {
	User       string
	Pass       string
	Host       string
	Connection *mgo.Session
	logger     lager.Logger
}

func New(username string, password string, host string, logger lager.Logger) (MongoProvisionerInterface, error) {
	var err error
	provisioner := MongoProvisioner{User: username, Pass: password, Host: host, logger: logger}
	provisioner.Connection, err = mgo.Dial(host)
	if err != nil {
		fmt.Println(err)
		provisioner.logger.Error("Error creating new mongo provisioner", err)
		return nil, err
	}
	return &provisioner, nil
}

func (e *MongoProvisioner) Close() {
	e.Connection.Close()
}

func (e *MongoProvisioner) IsDatabaseCreated(databaseName string) (bool, error) {
	databases, err := e.Connection.DatabaseNames()
	if err != nil {
		return false, err
	}

	for _, db := range databases {
		if db == databaseName {
			return true, nil
		}
	}

	return false, nil
}

func (e *MongoProvisioner) IsUserCreated(databaseName string, userName string) (bool, error) {
	userDB := e.Connection.DB(databaseName)
	result := bson.M{}
	err := userDB.Run(bson.M{"usersInfo": userName}, &result)
	if err != nil {
		return false, err
	}
	userInfo := result["users"].([]interface{})
	if len(userInfo) > 0 {
		return true, nil
	}
	return false, nil
}

func (e *MongoProvisioner) CreateDatabase(databaseName string) error {
	//this should create the db with empty users collection
	coll := e.Connection.DB(databaseName).C("sample")
	coll.Insert(bson.M{"a": 1, "b": 2})

	result := struct{ A, B int }{}

	err := coll.Find(bson.M{"a": 1}).One(&result)
	err = coll.DropCollection()
	if err != nil {
		return err
	}

	return nil
}

func (e *MongoProvisioner) DeleteDatabase(databaseName string) error {
	err := e.Connection.DB(databaseName).DropDatabase()
	if err != nil {
		return err
	}
	return nil
}

func (e *MongoProvisioner) CreateUser(databaseName string, username string, password string) error {
	userDB := e.Connection.DB(databaseName)

	err := userDB.AddUser(username, password, false)
	if err != nil {
		return err
	}
	return nil
}

func (e *MongoProvisioner) DeleteUser(databaseName string, username string) error {
	err := e.Connection.DB(databaseName).RemoveUser(username)
	if err != nil {
		return err
	}
	return nil
}
