package csm

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-openapi/runtime"
	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/hpcloud/cf-usb/lib/csm/client/connection"
	"github.com/hpcloud/cf-usb/lib/csm/client/workspace"
	"github.com/hpcloud/cf-usb/lib/csm/models"
	"github.com/pivotal-golang/lager"
)

type csmClient struct {
	logger           lager.Logger
	workspaceCient   *workspace.Client
	connectionClient *connection.Client
	authInfoWriter   runtime.ClientAuthInfoWriter
}

func NewCSMClient(logger lager.Logger, csmUrl url.URL, authToken string) CSMInterface {
	csm := csmClient{}
	transport := runtimeClient.New(csmUrl.Host, "/", []string{csmUrl.Scheme})
	csm.workspaceCient = workspace.New(transport, strfmt.Default)
	csm.connectionClient = connection.New(transport, strfmt.Default)
	csm.logger = logger
	csm.authInfoWriter = runtimeClient.APIKeyAuth("x-csm-token", "header", authToken)
	return &csm
}

func (csm *csmClient) CreateWorkspace(workspaceID string) error {
	csm.logger.Info("csm-create-workspace", lager.Data{"workspaceID": workspaceID})
	request := models.ServiceManagerWorkspaceCreateRequest{
		WorkspaceID: &workspaceID,
	}
	params := workspace.CreateWorkspaceParams{}
	params.CreateWorkspaceRequest = &request
	response, err := csm.workspaceCient.CreateWorkspace(&params, csm.authInfoWriter)

	if err != nil {
		return err
	}

	csm.logger.Info("csm-create-workspace", lager.Data{"response": response.Error()})

	//TODO: This is not working in CSM
	//	responseError := strings.TrimSpace(response.Error())
	//	if responseError != "" {
	//		logger.Info("csm-create-workspace", lager.Data{"error": "2"})
	//		return errors.New(responseError)
	//	}
	status := strings.TrimSpace(*response.Payload.Status)

	if status != "successful" {
		return errors.New(fmt.Sprintf("Error making the request. Extension returned status %s. Details: %v", status, response.Payload.Details))
	}

	return nil

}
func (csm *csmClient) WorkspaceExists(workspaceID string) (bool, error) {
	csm.logger.Info("csm-workspace-exists", lager.Data{"workspaceID": workspaceID})

	params := workspace.GetWorkspaceParams{}
	params.WorkspaceID = workspaceID
	response, err := csm.workspaceCient.GetWorkspace(&params, csm.authInfoWriter)

	if err != nil {
		return false, err
	}

	//TODO: This is not working in CSM
	//	responseError := strings.TrimSpace(response.Error())
	//	if responseError != "" {
	//		return false, errors.New(responseError)
	//	}
	csm.logger.Info("csm-workspace-exists", lager.Data{"response": response.Error()})

	status := strings.TrimSpace(*response.Payload.Status)

	//TODO: This is wrong and needs to be improved in the CSM server
	if status == "failed" {
		return false, nil
	}

	return true, nil

}
func (csm *csmClient) DeleteWorkspace(workspaceID string) error {
	csm.logger.Info("csm-delete-workspace", lager.Data{"workspaceID": workspaceID})
	params := workspace.DeleteWorkspaceParams{}
	params.WorkspaceID = workspaceID
	response, err := csm.workspaceCient.DeleteWorkspace(&params, csm.authInfoWriter)

	if err != nil {
		return err
	}

	csm.logger.Info("csm-delete-workspace", lager.Data{"response": response.Error()})
	//TODO: This is not working in CSM
	//	responseError := strings.TrimSpace(response.Error())
	//	if responseError != "" {
	//		return errors.New(responseError)
	//	}

	//TODO: in CSM this passes all the time, it does not take into consideration if a workspace exists

	return nil

}
func (csm *csmClient) CreateConnection(workspaceID, connectionID string) (interface{}, error) {
	csm.logger.Info("csm-create-connection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})
	params := connection.CreateConnectionParams{}
	params.WorkspaceID = workspaceID

	request := models.ServiceManagerConnectionCreateRequest{
		ConnectionID: &connectionID,
	}

	params.ConnectionCreateRequest = &request
	response, err := csm.connectionClient.CreateConnection(&params, csm.authInfoWriter)
	if err != nil {
		return nil, err
	}

	csm.logger.Info("csm-create-connection", lager.Data{"response": response.Error()})
	//TODO: This is not working in CSM
	//	responseError := strings.TrimSpace(response.Error())
	//	if responseError != "" {
	//		return errors.New(responseError)
	//	}
	status := strings.TrimSpace(*response.Payload.Status)

	if status != "successful" {
		return nil, errors.New(fmt.Sprintf("Error making the request. Extension returned status %s. Details: %v", status, response.Payload.Details))
	}

	return response.Payload.Details, err

}
func (csm *csmClient) ConnectionExists(workspaceID, connectionID string) (bool, error) {
	csm.logger.Info("csm-connection-exists", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})
	params := connection.GetConnectionParams{
		WorkspaceID:  workspaceID,
		ConnectionID: connectionID,
	}

	response, err := csm.connectionClient.GetConnection(&params, csm.authInfoWriter)
	if err != nil {
		return false, err
	}

	csm.logger.Info("csm-create-connection", lager.Data{"response": response.Error()})
	//TODO: This is not working in CSM
	//	responseError := strings.TrimSpace(response.Error())
	//	if responseError != "" {
	//		return errors.New(responseError)
	//	}

	status := strings.TrimSpace(*response.Payload.Status)

	//TODO: This is wrong and needs to be improved in the CSM server.
	//Currently is the only way to determine if a connection does not exist
	if status == "failed" {
		return false, nil
	}

	return true, nil

}
func (csm *csmClient) DeleteConnection(workspaceID, connectionID string) error {
	csm.logger.Info("csm-delete-connection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})
	params := connection.DeleteConnectionParams{
		WorkspaceID:  workspaceID,
		ConnectionID: connectionID,
	}

	response, err := csm.connectionClient.DeleteConnection(&params, csm.authInfoWriter)
	if err != nil {
		return err
	}

	csm.logger.Info("csm-delete-connection", lager.Data{"response": response.Error()})
	//TODO: This is not working in CSM
	//	responseError := strings.TrimSpace(response.Error())
	//	if responseError != "" {
	//		return errors.New(responseError)
	//	}

	//TODO: in CSM this passes all the time, it does not take into consideration if a connection exists
	return nil

}
