package config

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // mysql driver needed for database connection
	"github.com/hpcloud/cf-usb/lib/brokermodel"
	"github.com/hpcloud/cf-usb/lib/config/mysql/migrations"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/mysql"
	bindata "github.com/mattes/migrate/source/go-bindata"
)

type mysqlConfig struct {
	db         *sql.DB
	dbName     string
	configPath string
}

type generalConfig struct {
	key       string
	value     string
	component string
}

//NewMysqlConfig generates and returns a new mysql config provider
func NewMysqlConfig(address, username, password, database, configPath string) (Provider, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/?multiStatements=true", username, password, address))
	if err != nil {
		return nil, err
	}
	return &mysqlConfig{
		db:         db,
		dbName:     database,
		configPath: configPath,
	}, nil
}

func (c *mysqlConfig) InitializeConfiguration() error {
	if c.db == nil {
		return fmt.Errorf("Database connection not opened")
	}

	// Create the database if it doesn't exist
	if _, err := c.db.Exec("CREATE SCHEMA IF NOT EXISTS ? DEFAULT CHARACTER SET utf8", c.dbName); err != nil {
		return err
	}
	if _, err := c.db.Exec("USE ?", c.dbName); err != nil {
		return err
	}

	target, err := mysql.WithInstance(c.db, &mysql.Config{DatabaseName: c.dbName})
	if err != nil {
		return err
	}
	assets := bindata.Resource(migrations.AssetNames(), func(name string) ([]byte, error) {
		return migrations.Asset(name)
	})
	source, err := bindata.WithInstance(assets)
	if err != nil {
		return err
	}
	migration, err := migrate.NewWithInstance("go-bindata", source, "mysql", target)
	if err != nil {
		return err
	}

	return migration.Up()
}

func (c *mysqlConfig) LoadConfiguration() (*Config, error) {
	var configuration *Config
	var err error

	if c.configPath != "" {
		// Load it via the file provider
		configuration, err = NewFileConfig(c.configPath).LoadConfiguration()
		if err != nil {
			return nil, err
		}
		if configuration.ManagementAPI == nil {
			configuration.ManagementAPI = &ManagementAPI{}
		}
	}

	if configuration == nil {
		configuration = &Config{}
	}
	if configuration.ManagementAPI == nil {
		configuration.ManagementAPI = &ManagementAPI{}
	}

	configuration.Instances = make(map[string]Instance)

	instances, err := c.db.Query("SELECT Guid FROM Instances")
	if err != nil {
		return nil, err
	}

	defer instances.Close()

	for instances.Next() {
		var instanceGUID string
		err := instances.Scan(&instanceGUID)
		if err != nil {
			return nil, err
		}
		instance, err := c.LoadDriverInstance(instanceGUID)
		if err != nil {
			return nil, err
		}
		configuration.Instances[instanceGUID] = *instance
	}

	return configuration, nil
}

func (c *mysqlConfig) SaveConfiguration(config Config, overwrite bool) error {
	if overwrite == true {
		transaction, err := c.db.Begin()
		if err != nil {
			return err
		}

		transaction.Exec("DELETE FROM Config")

		transaction.Exec("INSERT INTO Config VALUES(?,?,?)", "API", config.APIVersion, "API_VERSION")
		transaction.Exec("INSERT INTO Config VALUES(?,?,?)", "EXTERNAL_URL", config.BrokerAPI.ExternalURL, "BROKER_API")
		transaction.Exec("INSERT INTO Config VALUES(?,?,?)", "LISTEN", config.BrokerAPI.Listen, "BROKER_API")
		transaction.Exec("INSERT INTO Config VALUES(?,?,?)", "REQUIRE_TLS", config.BrokerAPI.RequireTLS, "BROKER_API")
		transaction.Exec("INSERT INTO Config VALUES(?,?,?)", "SERVER_CERT_FILE", config.BrokerAPI.ServerCertFile, "BROKER_API")
		transaction.Exec("INSERT INTO Config VALUES(?,?,?)", "SERVER_KEY_FILE", config.BrokerAPI.ServerKeyFile, "BROKER_API")
		transaction.Exec("INSERT INTO Config VALUES(?,?,?)", "USERNAME", config.BrokerAPI.Credentials.Username, "BROKER_CREDENTIALS")
		transaction.Exec("INSERT INTO Config VALUES(?,?,?)", "PASSWORD", config.BrokerAPI.Credentials.Password, "BROKER_CREDENTIALS")
		transaction.Exec("INSERT INTO Config VALUES(?,?,?)", "BROKER_NAME", config.ManagementAPI.BrokerName, "MANAGEMENT_API")
		transaction.Exec("INSERT INTO Config VALUES(?,?,?)", "DEV_MODE", config.ManagementAPI.DevMode, "MANAGEMENT_API")
		transaction.Exec("INSERT INTO Config VALUES(?,?,?)", "LISTEN", config.ManagementAPI.Listen, "MANAGEMENT_API")
		transaction.Exec("INSERT INTO Config VALUES(?,?,?)", "UAA_CLIENT", config.ManagementAPI.UaaClient, "MANAGEMENT_API")
		transaction.Exec("INSERT INTO Config VALUES(?,?,?)", "UAA_SECRET", config.ManagementAPI.UaaSecret, "MANAGEMENT_API")
		transaction.Exec("INSERT INTO Config VALUES(?,?,?)", "API", config.ManagementAPI.CloudController.API, "CLOUD_CONTROLLER")
		transaction.Exec("INSERT INTO Config VALUES(?,?,?)", "SKIP_TLS_VALIDATION", config.ManagementAPI.CloudController.SkipTLSValidation, "CLOUD_CONTROLLER")

		if config.RoutesRegister != nil {
			transaction.Exec("INSERT INTO Config VALUES(?,?,?)", "BROKER_API_HOST", config.RoutesRegister.BrokerAPIHost, "ROUTES_REGISTER")
			transaction.Exec("INSERT INTO Config VALUES(?,?,?)", "MANAGEMENT_API_HOST", config.RoutesRegister.ManagmentAPIHost, "ROUTES_REGISTER")

			for _, member := range config.RoutesRegister.NatsMembers {
				transaction.Exec("INSERT INTO Config VALUES(?,?,?)", "NATS_MEMEBER", member, "ROUTES_REGISTER")
			}
		}

		transaction.Exec("INSERT INTO Config VALUES(?,?,?)", "AUTHENTICATION", string(*config.ManagementAPI.Authentication), "MANAGEMENT_API")

		err = transaction.Commit()
		if err != nil {
			transaction.Rollback()
			return err
		}
	}

	for instanceID, instance := range config.Instances {
		err := c.SetInstance(instanceID, instance)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *mysqlConfig) LoadDriverInstance(driverInstanceID string) (*Instance, error) {
	var driver Instance
	driver.Dials = make(map[string]Dial)

	result, err := c.db.Query("SELECT * FROM Instances WHERE Guid=?", driverInstanceID)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var instanceGUID string

	for result.Next() {
		if err := result.Scan(&instanceGUID, &driver.Name, &driver.TargetURL, &driver.AuthenticationKey, &driver.CaCert, &driver.SkipSsl); err != nil {
			return nil, err
		}
	}

	dials, err := c.db.Query("SELECT * FROM Dials WHERE Instances_Guid=?", instanceGUID)
	if err != nil {
		return nil, err
	}
	defer dials.Close()

	for dials.Next() {
		var dial Dial
		var configuration []byte
		var dialGUID string
		var planGUID string

		if err := dials.Scan(&dialGUID, &configuration, &planGUID, &instanceGUID); err != nil {
			return nil, err
		}
		rawConfig := json.RawMessage(configuration)
		dial.Configuration = &rawConfig

		planRow := c.db.QueryRow("SELECT * FROM Plans WHERE Guid=?", planGUID)

		var plan brokermodel.Plan
		var metadata brokermodel.PlanMetadata
		var meta []byte
		if err := planRow.Scan(&planGUID, &plan.Name, &plan.Description, &plan.Free, &meta); err != nil {
			return nil, err
		}
		err := json.Unmarshal(meta, &metadata)
		if err != nil {
			return nil, err
		}
		plan.Metadata = &metadata
		plan.ID = planGUID
		dial.Plan = plan

		driver.Dials[dialGUID] = dial
	}

	serviceRow := c.db.QueryRow("SELECT * FROM Services WHERE Instances_Guid=?", instanceGUID)

	var service brokermodel.CatalogService
	var dashboard []byte
	var metaService []byte
	var tags []byte
	var requires []byte
	if err := serviceRow.Scan(&service.ID, &service.Bindable, &dashboard, &service.Description, &metaService, &service.Name, &service.PlanUpdateable, &tags, &instanceGUID, &requires); err != nil {
		return nil, err
	}

	err = json.Unmarshal(metaService, &service.Metadata)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dashboard, &service.DashboardClient)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(tags, &service.Tags)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(requires, &service.Requires)
	if err != nil {
		return nil, err
	}

	driver.Service = service
	return &driver, nil
}

func (c *mysqlConfig) GetUaaAuthConfig() (*UaaAuth, error) {
	config, err := c.LoadConfiguration()
	if err != nil {
		return nil, err
	}

	conf := (*json.RawMessage)(config.ManagementAPI.Authentication)
	fmt.Println(string(*conf))
	uaa := Uaa{}
	err = json.Unmarshal(*conf, &uaa)
	if err != nil {
		return nil, err
	}
	return &uaa.UaaAuth, nil
}

func (c *mysqlConfig) SetInstance(instanceID string, instance Instance) error {
	_, err := c.db.Exec("INSERT INTO Instances VALUES(?, ?, ?,?,?,?)", instanceID, instance.Name, instance.TargetURL, instance.AuthenticationKey, instance.CaCert, instance.SkipSsl)
	if err != nil {
		return err
	}

	if len(instance.Dials) > 0 {
		for dialKey, dialInfo := range instance.Dials {
			err = c.SetDial(instanceID, dialKey, dialInfo)
			if err != nil {
				return err
			}
		}
	}

	if instance.Service.Name != "" {
		err = c.SetService(instanceID, instance.Service)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *mysqlConfig) GetInstance(instanceID string) (*Instance, string, error) {
	instanceRow := c.db.QueryRow("SELECT * FROM Instances WHERE Guid=?", instanceID)
	var instance Instance
	var instanceGUID string

	err := instanceRow.Scan(&instanceGUID, &instance.Name, &instance.TargetURL, &instance.AuthenticationKey, &instance.CaCert, &instance.SkipSsl)
	if err != nil {
		return nil, "", err
	}
	return &instance, instanceGUID, nil
}

func (c *mysqlConfig) DeleteInstance(instanceID string) error {

	dials, err := c.db.Query("SELECT * FROM Dials WHERE Instances_Guid=?", instanceID)
	defer dials.Close()

	transaction, err := c.db.Begin()
	if err != nil {
		return err
	}

	for dials.Next() {
		var configuration []byte
		var dialGUID string
		var planGUID string
		var instanceGUID string

		if err := dials.Scan(&dialGUID, &configuration, &planGUID, &instanceGUID); err != nil {
			return err
		}
		_, err = transaction.Exec("DELETE FROM Dials WHERE Guid=?", dialGUID)

		_, err = transaction.Exec("DELETE FROM Plans WHERE Guid=?", planGUID)
	}

	_, err = transaction.Exec("DELETE FROM Services WHERE Instances_Guid=?", instanceID)

	_, err = transaction.Exec("DELETE FROM Instances WHERE Guid=?", instanceID)
	if err != nil {
		err = transaction.Rollback()
		return err
	}
	err = transaction.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (c *mysqlConfig) SetService(instanceID string, service brokermodel.CatalogService) error {
	dashboard, err := json.Marshal(service.DashboardClient)
	if err != nil {
		return err
	}
	metadata, err := json.Marshal(service.Metadata)
	if err != nil {
		return err
	}
	tags, err := json.Marshal(service.Tags)
	if err != nil {
		return err
	}
	requires, err := json.Marshal(service.Requires)
	if err != nil {
		return err
	}
	_, err = c.db.Exec("INSERT INTO Services VALUES(?,?,?,?,?,?,?,?,?,?)", service.ID, service.Bindable, dashboard, service.Description, metadata, service.Name, service.PlanUpdateable, tags, instanceID, requires)
	if err != nil {
		return err
	}
	return nil
}

func (c *mysqlConfig) GetService(serviceID string) (*brokermodel.CatalogService, string, error) {

	serviceRow := c.db.QueryRow("SELECT * FROM Services WHERE Guid=?", serviceID)
	var service brokermodel.CatalogService
	var dash []byte
	var meta []byte
	var tags []byte
	var requires []byte
	var instanceID string
	err := serviceRow.Scan(&service.ID, &service.Bindable, &dash, &service.Description, &meta, &service.Name, &service.PlanUpdateable, &tags, &instanceID, &requires)

	var dashboard brokermodel.DashboardClient
	err = json.Unmarshal(dash, &dashboard)
	if err != nil {
		return nil, "", err
	}
	service.DashboardClient = &dashboard

	err = json.Unmarshal(meta, &service.Metadata)
	if err != nil {
		return nil, "", err
	}
	err = json.Unmarshal(tags, &service.Tags)
	if err != nil {
		return nil, "", err
	}

	err = json.Unmarshal(requires, &service.Requires)
	if err != nil {
		return nil, "", err
	}

	return &service, instanceID, nil
}

func (c *mysqlConfig) DeleteService(instanceID string) error {
	var serviceGUID string
	if err := c.db.QueryRow("SELECT Services_Guid FROM Instances WHERE Guid=?", instanceID).Scan(&serviceGUID); err != nil {
		return err
	}

	_, err := c.db.Exec("DELETE FROM Services WHERE Guid=?", serviceGUID)
	if err != nil {
		return err
	}

	return nil
}

func (c *mysqlConfig) SetDial(instanceID string, dialID string, dial Dial) error {
	configuration, err := json.Marshal(dial.Configuration)
	if err != nil {
		return err
	}

	meta, err := json.Marshal(dial.Plan.Metadata)
	if err != nil {
		return err
	}

	transaction, err := c.db.Begin()
	if err != nil {
		return err
	}
	_, err = transaction.Exec("INSERT INTO Plans VALUES(?, ?, ?,?, ?)", dial.Plan.ID, dial.Plan.Name, dial.Plan.Description, dial.Plan.Free, meta)

	_, err = transaction.Exec("INSERT INTO Dials VALUES(?, ?, ?,?)", dialID, configuration, dial.Plan.ID, instanceID)

	if err != nil {
		err = transaction.Rollback()
		return err
	}
	transaction.Commit()

	return nil
}

func (c *mysqlConfig) GetDial(dialID string) (*Dial, string, error) {

	dialRow := c.db.QueryRow("SELECT * FROM Dials WHERE Guid=?", dialID)

	var dialGUID string
	var conf []byte
	var planGUID string
	var instanceGUID string
	err := dialRow.Scan(&dialGUID, &conf, &planGUID, &instanceGUID)

	var config json.RawMessage

	err = json.Unmarshal(conf, &config)
	if err != nil {
		return nil, "", err
	}

	plan, _, _, err := c.GetPlan(planGUID)
	if err != nil {
		return nil, "", err
	}
	var result Dial
	result.Configuration = &config
	result.Plan = *plan
	return &result, instanceGUID, nil
}

func (c *mysqlConfig) DeleteDial(dialID string) error {

	dial, instanceID, err := c.GetDial(dialID)
	if err != nil {
		return err
	}

	transaction, err := c.db.Begin()
	if err != nil {
		return err
	}
	_, err = transaction.Exec("DELETE FROM Plans WHERE Guid=?", dial.Plan.ID)

	_, err = transaction.Exec("DELETE FROM Dials WHERE Instances_Guid=?", instanceID)
	if err != nil {
		err = transaction.Rollback()
		return err
	}
	err = transaction.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (c *mysqlConfig) InstanceNameExists(driverInstanceName string) (bool, error) {
	result := c.db.QueryRow("SELECT EXISTS(SELECT * FROM Instances WHERE Name=?)", driverInstanceName)
	var exists bool

	if err := result.Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (c *mysqlConfig) GetPlan(planid string) (*brokermodel.Plan, string, string, error) {

	planRow := c.db.QueryRow("SELECT * FROM Plans WHERE Guid=?", planid)

	var plan brokermodel.Plan
	var meta []byte
	var planGUID string

	err := planRow.Scan(&planGUID, &plan.Name, &plan.Description, &plan.Free, &meta)
	if err != nil {
		return nil, "", "", err
	}
	var metadata brokermodel.PlanMetadata
	err = json.Unmarshal(meta, &metadata)
	if err != nil {
		return nil, "", "", err
	}
	plan.Metadata = &metadata
	plan.ID = planGUID

	var dialID string
	var instanceID string

	err = c.db.QueryRow("SELECT Guid,Instances_Guid FROM Dials WHERE Plans_Guid=?", planid).Scan(&dialID, &instanceID)
	if err != nil {
		return nil, "", "", err
	}
	return &plan, dialID, instanceID, nil
}
