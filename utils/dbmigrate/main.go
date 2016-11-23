package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/consul/api"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/config/consul"
)

//Migration tool to allow data migration from consul kv store to mysql database
//Warning: this tools empties existing mysql configuration table before migration

func main() {
	consulAddress := os.Getenv("CONSUL_ADDRESS")

	if consulAddress == "" {
		fmt.Println("CONSUL_ADDRESS must be set")
		os.Exit(0)
	}

	consulDatacenter := os.Getenv("CONSUL_DATACENTER")
	consulUser := os.Getenv("CONSUL_USERNAME")
	consulPass := os.Getenv("CONSUL_PASSWORD")
	consulSchema := os.Getenv("CONSUL_SCHEMA")
	consulToken := os.Getenv("CONSUL_TOKEN")

	var consulConfig api.Config
	consulConfig.Address = consulAddress
	consulConfig.Datacenter = consulDatacenter

	var auth api.HttpBasicAuth
	auth.Username = consulUser
	auth.Password = consulPass

	consulConfig.HttpAuth = &auth
	consulConfig.Scheme = consulSchema

	consulConfig.Token = consulToken

	provisioner, err := consul.New(&consulConfig)
	if err != nil {
		fmt.Println("consul config provider", err)
		os.Exit(1)
	}

	configuraiton := config.NewConsulConfig(provisioner)
	configData, err := configuraiton.LoadConfiguration()
	if err != nil {
		fmt.Println("load consul configuration ", err)
		os.Exit(1)
	}

	mysqlAddress := os.Getenv("MYSQL_ADDRESS")
	mysqlDB := os.Getenv("MYSQL_DB")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPass := os.Getenv("MYSQL_PASSWORD")

	if mysqlAddress == "" || mysqlDB == "" || mysqlUser == "" || mysqlPass == "" {
		fmt.Println("MYSQL configration environment variables must be set (MYSQL_ADDRESS,MYSQL_DB,MYSQL_USER,MYSQL_PASSWORD")
		os.Exit(1)
	}

	mysqlConfiguration, err := config.NewMysqlConfig(mysqlAddress, mysqlUser, mysqlPass, mysqlDB)
	if err != nil {
		fmt.Println("mysql config provider error", err)
		os.Exit(1)
	}

	overwrite := false
	overwriteConfiguration := os.Getenv("OVERWRITE_CONFIG")
	if overwriteConfiguration != "" {
		overwrite, err = strconv.ParseBool(overwriteConfiguration)
		if err != nil {
			fmt.Println("Error parsing overwrite config environment variable", err)
			os.Exit(1)
		}
	}

	existing, _ := mysqlConfiguration.LoadConfiguration()

	if existing != nil {
		if len(existing.Instances) > 0 {
			for instanceID, _ := range existing.Instances {
				err = mysqlConfiguration.DeleteInstance(instanceID)
				if err != nil {
					fmt.Println("Error cleaning up mysql instances", err)
					os.Exit(1)
				}
			}
		}
	}

	err = mysqlConfiguration.SaveConfiguration(*configData, overwrite)
	if err != nil {
		fmt.Println("save mysql configuration ", err)
		os.Exit(1)
	}

}
