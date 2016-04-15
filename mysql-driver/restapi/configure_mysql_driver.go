package restapi

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	errors "github.com/go-swagger/go-swagger/errors"
	httpkit "github.com/go-swagger/go-swagger/httpkit"
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"

	"database/sql" //added
	"mysql-driver/models"
	"mysql-driver/restapi/operations"
	"mysql-driver/restapi/operations/connection"
	"mysql-driver/restapi/operations/workspace"

	_ "github.com/go-sql-driver/mysql" //added
)

/********************/
/*added section*/

//const host string = "mydb.com1dc8kxmtn.eu-central-1.rds.amazonaws.com"
//const port string = "3306"
//const user string = "mitza"
//const pass string = "password"

var host string
var port string
var user string
var pass string

var web_port string

func getDBNameFromId(id string) string {
	dbName := "d" + strings.Replace(id, "-", "", -1)
	dbName = strings.Replace(dbName, ";", "", -1)

	return dbName
}

func getConnectionString() (string, error) {
	host = os.Getenv("MYSQL_HOST")
	port = os.Getenv("MYSQL_PORT")
	user = os.Getenv("MYSQL_USER")
	pass = os.Getenv("MYSQL_PASS")

	if host != "" && port != "" && user != "" && pass != "" {
		conn_string := user + ":" + pass + "@tcp(" + host + ":" + port + ")/"
		return conn_string, nil
	}

	return "", errors.NotFound("You should set the MYSQL_HOST, MYSQL_PORT, MYSQL_USER and MYSQL_PASS environment variables!")
}

func NewMySqlDb(response *models.ServiceManagerWorkspaceResponse, workspace_id string) error {
	id := workspace_id

	response.ProcessingType = "Default"

	if err := CreateDatabase(id); err != nil {
		fmt.Println(err.Error())
		response.Status = "failed"
		return err
	}

	response.Status = "successful"
	return nil
}

func DeleteMySqlDb(workspace_id string) error {

	if err := DeleteDatabase(workspace_id); err != nil {
		return err
	}

	return nil
}

func CheckMySqlDb(workspace_id string, response *models.ServiceManagerWorkspaceResponse) (bool, error) {
	dbExists, err := CheckIfDatabaseExists(workspace_id)

	if err != nil {
		return false, err
	} else {
		response.Status = "successful"
		response.ProcessingType = "Default"
		if !dbExists {
			details := make(map[string]interface{})
			details["database"] = "not found"
			response.Details = details
			return false, nil
		}
	}
	return true, nil
}

func BindInstance(workspace_id string, connection_id string, response *models.ServiceManagerConnectionResponse) error {

	id := workspace_id
	bind_id := connection_id

	response.ProcessingType = "Default"

	content, err := CreateUserForDatabase(id, bind_id)

	if err != nil {
		fmt.Println(err.Error())
		response.Status = "failed"
		return err
	}

	details := make(map[string]interface{})
	details[content["username"]] = content

	response.Status = "successful"
	response.Details = details
	return nil
}

func DeleteBind(workspace_id string, connection_id string) error {
	return RemoveGrant(workspace_id, connection_id)
}

func CheckBind(workspace_id string, connection_id string, response *models.ServiceManagerConnectionResponse) (bool, error) {
	response.ProcessingType = "Default"

	rez, err := CheckIfUserHasDBRights(workspace_id, connection_id)

	response.ProcessingType = "Default"

	if err != nil {

		response.Status = "failed"
		return false, err
	}

	response.Status = "successful"

	details := make(map[string]interface{})
	details[rez["username"]] = rez

	response.Details = details

	if rez["has_db_access"] == "false" {
		return false, nil
	}

	return true, nil
}

func CheckIfUserHasDBRights(database string, user string) (map[string]string, error) {
	conn_string, err := getConnectionString()
	if err != nil {
		return nil, err
	}
	rez := make(map[string]string)

	dbName := getDBNameFromId(database)

	username, err := getMD5Hash(user)
	if err != nil {
		return rez, err
	}
	if len(username) > 16 {
		username = username[:16]
	}
	db, err := sql.Open("mysql", conn_string+dbName)

	if err != nil {
		return rez, err
	}

	defer db.Close()

	rows, err := db.Query("Select count(*) from mysql.db where user=? and db=?", username, dbName)

	if err != nil {
		return rez, err
	}

	var count int

	for rows.Next() {
		rows.Scan(&count)
	}

	var has_db_access string
	if count == 0 {
		has_db_access = "false"
	} else {
		has_db_access = "true"
	}
	rez["hostname"] = host
	rez["host"] = host
	rez["user"] = username
	rez["port"] = port
	rez["username"] = username
	rez["database"] = dbName
	rez["has_db_access"] = has_db_access

	return rez, nil

}

//connect to mysql and create an user that will have select, update, delete,
//create etc. rights
func CreateUserForDatabase(id string, bind_id string) (map[string]string, error) {

	conn_string, err := getConnectionString()
	if err != nil {
		return nil, err
	}
	rez := make(map[string]string)

	dbName := getDBNameFromId(id)

	username, err := getMD5Hash(bind_id)
	if err != nil {
		return rez, err
	}
	if len(username) > 16 {
		username = username[:16]
	}

	password, _ := secureRandomString(32)

	db, err := sql.Open("mysql", conn_string+dbName)

	if err != nil {
		return rez, err
	}

	defer db.Close()
	sql := fmt.Sprintf("grant select,update,delete, create, execute, show view, alter, alter routine "+
		", create routine,create view, index, drop,references, create temporary tables, lock tables on %s"+
		".* to '%s'@'%%' identified by '%s'", dbName, username, password)

	_, err = db.Query(sql)

	if err != nil {
		return rez, err
	}

	rez["hostname"] = host
	rez["host"] = host
	rez["user"] = bind_id
	rez["port"] = port
	rez["username"] = bind_id
	rez["password"] = password
	rez["database"] = dbName

	return rez, err
}

func getMD5Hash(text string) (string, error) {
	hasher := md5.New()
	hasher.Write([]byte(text))
	generated := hex.EncodeToString(hasher.Sum(nil))

	reg := regexp.MustCompile("[^A-Za-z0-9]+")

	return reg.ReplaceAllString(generated, ""), nil
}

func secureRandomString(bytesOfEntpry int) (string, error) {
	rb := make([]byte, bytesOfEntpry)
	_, err := rand.Read(rb)

	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(rb), nil
}

//deletes a database from mysql
func DeleteDatabase(id string) error {
	conn_string, err := getConnectionString()

	if err != nil {
		return err
	}

	dbName := getDBNameFromId(id)

	db, err := sql.Open("mysql", conn_string)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Query("Drop database " + dbName)

	if err != nil {
		return err
	}
	return nil
}

//creates a database in mysql
func CreateDatabase(id string) error {

	conn_string, err := getConnectionString()

	if err != nil {
		return err
	}

	dbName := getDBNameFromId(id)

	db, err := sql.Open("mysql", conn_string)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Query("Create database " + dbName)

	if err != nil {
		return err
	}
	return nil
}

func CheckIfDatabaseExists(id string) (bool, error) {
	conn_string, err := getConnectionString()

	if err != nil {
		return false, err
	}

	dbName := getDBNameFromId(id)

	db, err := sql.Open("mysql", conn_string)
	if err != nil {
		return false, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '" + dbName + "'")

	if err != nil {
		return false, err
	}

	iCnt := 0

	var row string

	for rows.Next() {
		iCnt++
		rows.Scan(&row)
	}

	if iCnt == 1 {
		return true, nil
	}

	return false, nil
}

func RemoveGrant(id string, bind_id string) error {
	conn_string, err := getConnectionString()
	if err != nil {
		return err
	}
	//rez := make(map[string]string)

	dbName := getDBNameFromId(id)

	username, err := getMD5Hash(bind_id)
	if err != nil {
		return err
	}
	if len(username) > 16 {
		username = username[:16]
	}

	//password, _ := secureRandomString(32)

	db, err := sql.Open("mysql", conn_string+dbName)

	if err != nil {
		return err
	}
	defer db.Close()

	uiodb, err := UserInOtherDBs(username, db)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if !uiodb {
		_, err = db.Query("drop user '" + username + "'")

		return nil
	}
	if err != nil {
		return err
	}

	_, err = db.Query("REVOKE all privileges on " + dbName +
		".* from '" + username + "'@'%'")

	if err != nil {
		return err
	}
	return nil
}

func UserInOtherDBs(bind_id string, db *sql.DB) (bool, error) {

	rows, err := db.Query("Show grants for " + bind_id)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "Error 1141:") {
			return false, err
		} else {
			return false, nil //no grant - so no user
		}
	}

	var row string

	iCnt := 0

	for rows.Next() {
		iCnt++
		rows.Scan(&row)
	}

	return iCnt > 2, nil

}

/*finish added section*/
/*************************/

// This file is safe to edit. Once it exists it will not be overwritten

func configureFlags(api *operations.MysqlDriverAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.MysqlDriverAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.ConnectionCreateConnectionHandler = connection.CreateConnectionHandlerFunc(func(params connection.CreateConnectionParams) middleware.Responder {

		r := models.ServiceManagerConnectionResponse{
			Details: params.ConnectionCreateRequest.Details,
			Status:  "none",
		}

		response := &r

		BindInstance(params.WorkspaceID, params.ConnectionCreateRequest.ConnectionID, response)

		return connection.NewCreateConnectionCreated().WithPayload(response)
		//return middleware.NotImplemented("operation connection.CreateConnection has not yet been implemented")
	})
	api.WorkspaceCreateWorkspaceHandler = workspace.CreateWorkspaceHandlerFunc(func(params workspace.CreateWorkspaceParams) middleware.Responder {
		r := models.ServiceManagerWorkspaceResponse{
			Details:        params.CreateWorkspaceRequest.Details,
			ProcessingType: "Default",
			Status:         "none"}
		response := &r

		err := NewMySqlDb(&r, params.CreateWorkspaceRequest.WorkspaceID)

		if err != nil {
			var errCode int64
			errCode = 500 //http code 500: Internal Server error
			errPayload := models.Error{&errCode, err.Error()}
			return workspace.NewDeleteWorkspaceDefault(500).WithPayload(&errPayload)
		}

		return workspace.NewCreateWorkspaceCreated().WithPayload(response)
	})
	api.ConnectionDeleteConnectionHandler = connection.DeleteConnectionHandlerFunc(func(params connection.DeleteConnectionParams) middleware.Responder {
		err := DeleteBind(params.WorkspaceID, params.ConnectionID)
		if err != nil {
			var errCode int64
			errCode = 500
			errPayload := models.Error{&errCode, err.Error()}
			return connection.NewDeleteConnectionDefault(500).WithPayload(&errPayload)
		}
		return connection.NewDeleteConnectionOK()
	})
	api.WorkspaceDeleteWorkspaceHandler = workspace.DeleteWorkspaceHandlerFunc(func(params workspace.DeleteWorkspaceParams) middleware.Responder {
		err := DeleteMySqlDb(params.WorkspaceID)
		if err != nil {
			var errCode int64
			if strings.HasPrefix(err.Error(), "Error 1008:") {
				errCode = 410 //http code 410: GONE
			} else {
				errCode = 500
			}
			errPayload := models.Error{&errCode, err.Error()}
			return workspace.NewDeleteWorkspaceDefault(int(*errPayload.Code)).WithPayload(&errPayload)
		}
		return workspace.NewDeleteWorkspaceOK()
	})
	api.ConnectionGetConnectionHandler = connection.GetConnectionHandlerFunc(func(params connection.GetConnectionParams) middleware.Responder {
		r := models.ServiceManagerConnectionResponse{
			ProcessingType: "None",
			Status:         "none",
		}
		has_db_access, err := CheckBind(params.WorkspaceID, params.ConnectionID, &r)

		if err != nil {
			var errCode int64
			errCode = 500
			errPayload := models.Error{&errCode, err.Error()}
			return connection.NewGetConnectionDefault(500).WithPayload(&errPayload)
		}
		if !has_db_access {
			var errCode int64
			errCode = 401
			errPayload := models.Error{&errCode, "The user has no access to this database"}
			return connection.NewGetConnectionDefault(410).WithPayload(&errPayload)
		}
		return connection.NewGetConnectionOK().WithPayload(&r)
	})
	api.WorkspaceGetWorkspaceHandler = workspace.GetWorkspaceHandlerFunc(func(params workspace.GetWorkspaceParams) middleware.Responder {
		r := models.ServiceManagerWorkspaceResponse{
			ProcessingType: "Default",
			Status:         "none",
		}
		exists, err := CheckMySqlDb(params.WorkspaceID, &r)
		if err != nil {
			var errCode int64
			errCode = 500
			errPayload := models.Error{&errCode, err.Error()}
			return workspace.NewGetWorkspaceDefault(500).WithPayload(&errPayload)
		} else {
			if exists {
				return workspace.NewGetWorkspaceOK().WithPayload(&r)
			} else {
				var errCode int64
				errCode = 410 //http code 410: GONE
				errPayload := models.Error{&errCode, "The database was not found"}
				return workspace.NewGetWorkspaceDefault(410).WithPayload(&errPayload) //gone
			}
		}

	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
