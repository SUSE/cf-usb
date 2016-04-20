package restapi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	errors "github.com/go-swagger/go-swagger/errors"
	httpkit "github.com/go-swagger/go-swagger/httpkit"
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/hpcloud/cf-usb/sidecar/clients/util"
	"github.com/hpcloud/cf-usb/sidecar/server/executablecaller"
	"github.com/hpcloud/cf-usb/sidecar/server/models"
	"github.com/hpcloud/cf-usb/sidecar/server/restapi/operations"
	"github.com/hpcloud/cf-usb/sidecar/server/restapi/operations/connection"
	"github.com/hpcloud/cf-usb/sidecar/server/restapi/operations/workspace"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureFlags(api *operations.BrokerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func GetCodeFromClientResponse(clientResponse string) int {
	lines := strings.Split(clientResponse, "\n")
	if len(lines) == 0 {
		return -2
	}
	fmt.Println(lines[0])
	rez, err := strconv.Atoi(strings.Trim(lines[0], "\r\n"))
	if err != nil {
		return -3
	}

	return rez
}

func configureAPI(api *operations.BrokerAPI) http.Handler {
	// configure the api here

	var caller executablecaller.IExecutableCaller
	caller = executablecaller.DefaultCaller{}

	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.ConnectionCreateConnectionHandler = connection.CreateConnectionHandlerFunc(func(params connection.CreateConnectionParams) middleware.Responder {
		r := models.ServiceManagerConnectionResponse{
			Details:        params.ConnectionCreateRequest.Details,
			Status:         "none",
			ProcessingType: "Default",
		}

		clientResponse, err := caller.CreateConnectionCaller(params.WorkspaceID, params.ConnectionCreateRequest.ConnectionID)

		if err == -1 {
			r.Status = "failed"
			var codeErr int64 = -1
			return connection.NewCreateConnectionDefault(500).WithPayload(&models.Error{Code: &codeErr, Message: fmt.Sprintf("Unknown error, timeout error or executable not present?")})
		}

		jsonResp := util.JsonResp{}
		errJson := json.Unmarshal([]byte(clientResponse), &jsonResp)

		var codeJson int64 = 500
		if errJson != nil {
			return connection.NewCreateConnectionDefault(500).WithPayload(&models.Error{Code: &codeJson, Message: "Json error"})
		}

		if err != 0 {
			r.Status = "failed"

			errModel := util.ErrorMsg{}
			errModel.Code = int(jsonResp.Payload.(map[string]interface{})["code"].(float64))
			if errModel.Code < 200 || errModel.Code > 500 {
				errModel.Code = 500
			}
			errModel.Message = jsonResp.Payload.(map[string]interface{})["message"].(string)
			return connection.NewCreateConnectionDefault(jsonResp.HttpCode).WithPayload(&models.Error{Code: &codeJson, Message: errModel.Message})

		}

		r.Status = "successful"
		details := jsonResp.Payload.(map[string]interface{})

		r.Details = details
		return connection.NewCreateConnectionCreated().WithPayload(&r)
	})
	api.WorkspaceCreateWorkspaceHandler = workspace.CreateWorkspaceHandlerFunc(func(params workspace.CreateWorkspaceParams) middleware.Responder {
		r := models.ServiceManagerWorkspaceResponse{
			Details:        params.CreateWorkspaceRequest.Details,
			Status:         "none",
			ProcessingType: "Default",
		}

		clientResponse, err := caller.CreateWorkspaceCaller(params.CreateWorkspaceRequest.WorkspaceID)

		if err == -1 {
			r.Status = "failed"
			var codeErr int64 = -1
			return workspace.NewCreateWorkspaceDefault(500).WithPayload(&models.Error{Code: &codeErr, Message: fmt.Sprintf("Unknown error, timeout error or executable not present?")})
		}

		jsonResp := util.JsonResp{}
		errJson := json.Unmarshal([]byte(clientResponse), &jsonResp)

		var codeJson int64 = 500
		if errJson != nil {
			return workspace.NewCreateWorkspaceDefault(500).WithPayload(&models.Error{Code: &codeJson, Message: "Json error"})
		}
		if err != 0 {
			r.Status = "failed"

			errModel := util.ErrorMsg{}
			errModel.Code = int(jsonResp.Payload.(map[string]interface{})["code"].(float64))
			if errModel.Code < 200 || errModel.Code > 500 {
				errModel.Code = 500
			}
			errModel.Message = jsonResp.Payload.(map[string]interface{})["message"].(string)
			return workspace.NewCreateWorkspaceDefault(jsonResp.HttpCode).WithPayload(&models.Error{Code: &codeJson, Message: errModel.Message})
		}

		r.Status = "successful"
		details := jsonResp.Payload.(map[string]interface{})

		r.Details = details
		return workspace.NewCreateWorkspaceCreated().WithPayload(&r)
	})
	api.ConnectionDeleteConnectionHandler = connection.DeleteConnectionHandlerFunc(func(params connection.DeleteConnectionParams) middleware.Responder {

		clientResponse, err := caller.DeleteConnectionCaller(params.ConnectionID)

		if err == -1 {
			var codeErr int64 = -1
			return connection.NewDeleteConnectionDefault(500).WithPayload(&models.Error{Code: &codeErr, Message: fmt.Sprintf("Unknown error, timeout error or executable not present?")})
		}
		jsonResp := util.JsonResp{}
		errJson := json.Unmarshal([]byte(clientResponse), &jsonResp)

		var codeJson int64 = 500
		if errJson != nil {
			return connection.NewDeleteConnectionDefault(500).WithPayload(&models.Error{Code: &codeJson, Message: "Json error"})
		}

		if err != 0 {
			errModel := util.ErrorMsg{}
			errModel.Code = int(jsonResp.Payload.(map[string]interface{})["code"].(float64))
			if errModel.Code < 200 || errModel.Code > 500 {
				errModel.Code = 500
			}
			errModel.Message = jsonResp.Payload.(map[string]interface{})["message"].(string)
			return connection.NewDeleteConnectionDefault(jsonResp.HttpCode).WithPayload(&models.Error{Code: &codeJson, Message: errModel.Message})

		}

		return connection.NewDeleteConnectionOK()
	})
	api.WorkspaceDeleteWorkspaceHandler = workspace.DeleteWorkspaceHandlerFunc(func(params workspace.DeleteWorkspaceParams) middleware.Responder {

		clientResponse, err := caller.DeleteWorkspaceCaller(params.WorkspaceID)

		if err == -1 {
			var codeErr int64 = -1
			return workspace.NewDeleteWorkspaceDefault(500).WithPayload(&models.Error{Code: &codeErr, Message: fmt.Sprintf("Unknown error, timeout error or executable not present?")})
		}
		jsonResp := util.JsonResp{}
		errJson := json.Unmarshal([]byte(clientResponse), &jsonResp)

		var codeJson int64 = 500
		if errJson != nil {
			return workspace.NewDeleteWorkspaceDefault(500).WithPayload(&models.Error{Code: &codeJson, Message: "Json error"})
		}

		if err != 0 {
			errModel := util.ErrorMsg{}
			errModel.Code = int(jsonResp.Payload.(map[string]interface{})["code"].(float64))
			if errModel.Code < 200 || errModel.Code > 500 {
				errModel.Code = 500
			}
			errModel.Message = jsonResp.Payload.(map[string]interface{})["message"].(string)
			return workspace.NewDeleteWorkspaceDefault(jsonResp.HttpCode).WithPayload(&models.Error{Code: &codeJson, Message: errModel.Message})

		}

		return workspace.NewDeleteWorkspaceOK()
	})
	api.ConnectionGetConnectionHandler = connection.GetConnectionHandlerFunc(func(params connection.GetConnectionParams) middleware.Responder {
		r := models.ServiceManagerConnectionResponse{
			Status:         "none",
			ProcessingType: "Default",
		}

		clientResponse, err := caller.GetConnectionCaller(params.WorkspaceID, params.ConnectionID)

		if err == -1 {
			r.Status = "failed"
			var codeErr int64 = -1
			return connection.NewGetConnectionDefault(500).WithPayload(&models.Error{Code: &codeErr, Message: fmt.Sprintf("Unknown error, timeout error or executable not present?")})
		}
		jsonResp := util.JsonResp{}
		errJson := json.Unmarshal([]byte(clientResponse), &jsonResp)

		var codeJson int64 = 500
		if errJson != nil {
			return connection.NewGetConnectionDefault(500).WithPayload(&models.Error{Code: &codeJson, Message: "Json error"})
		}

		if err != 0 {
			r.Status = "failed"

			errModel := util.ErrorMsg{}
			errModel.Code = int(jsonResp.Payload.(map[string]interface{})["code"].(float64))
			if errModel.Code < 200 || errModel.Code > 500 {
				errModel.Code = 500
			}
			errModel.Message = jsonResp.Payload.(map[string]interface{})["message"].(string)
			return connection.NewGetConnectionDefault(jsonResp.HttpCode).WithPayload(&models.Error{Code: &codeJson, Message: errModel.Message})

		}

		r.Status = "successful"
		details := jsonResp.Payload.(map[string]interface{})

		r.Details = details
		return connection.NewGetConnectionOK().WithPayload(&r)
	})
	api.WorkspaceGetWorkspaceHandler = workspace.GetWorkspaceHandlerFunc(func(params workspace.GetWorkspaceParams) middleware.Responder {
		r := models.ServiceManagerWorkspaceResponse{
			Status:         "none",
			ProcessingType: "Default",
		}

		clientResponse, err := caller.GetWorkspaceCaller(params.WorkspaceID)

		fmt.Println(clientResponse)

		if err == -1 {
			r.Status = "failed"
			var codeErr int64 = -1
			return workspace.NewGetWorkspaceDefault(500).WithPayload(&models.Error{Code: &codeErr, Message: fmt.Sprintf("Unknown error, timeout error or executable not present?")})
		}

		jsonResp := util.JsonResp{}
		errJson := json.Unmarshal([]byte(clientResponse), &jsonResp)

		var codeJson int64 = 500
		if errJson != nil {
			return workspace.NewGetWorkspaceDefault(500).WithPayload(&models.Error{Code: &codeJson, Message: "Json error"})
		}
		if err != 0 {
			r.Status = "failed"

			errModel := util.ErrorMsg{}
			errModel.Code = int(jsonResp.Payload.(map[string]interface{})["code"].(float64))
			if errModel.Code < 200 || errModel.Code > 500 {
				errModel.Code = 500
			}
			errModel.Message = jsonResp.Payload.(map[string]interface{})["message"].(string)
			return workspace.NewGetWorkspaceDefault(jsonResp.HttpCode).WithPayload(&models.Error{Code: &codeJson, Message: errModel.Message})
		}

		r.Status = "successful"
		details := jsonResp.Payload.(map[string]interface{})

		r.Details = details
		return workspace.NewGetWorkspaceOK().WithPayload(&r)
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
