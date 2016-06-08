package csm

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

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
	loggedIn         bool
}

//NewCSMClient instantiates a new csmClient
func NewCSMClient(logger lager.Logger) CSM {
	csm := csmClient{}
	csm.logger = logger
	csm.loggedIn = false
	return &csm
}

func (csm *csmClient) Login(targetEndpoint string, token string) error {
	csm.logger.Info("csm-login", lager.Data{"endpoint": targetEndpoint})
	target, err := url.Parse(targetEndpoint)
	if err != nil {
		return err
	}
	transport := runtimeClient.New(target.Host, "/", []string{target.Scheme})
	csm.workspaceCient = workspace.New(transport, strfmt.Default)
	csm.connectionClient = connection.New(transport, strfmt.Default)
	csm.authInfoWriter = runtimeClient.APIKeyAuth("x-csm-token", "header", token)
	csm.loggedIn = true
	return nil
}

func (csm *csmClient) CreateWorkspace(workspaceID string) error {
	if !csm.loggedIn {
		return errors.New("Not logged in")
	}
	csm.logger.Info("csm-create-workspace", lager.Data{"workspaceID": workspaceID})
	request := models.ServiceManagerWorkspaceCreateRequest{
		WorkspaceID: &workspaceID,
	}
	params := workspace.CreateWorkspaceParams{}
	params.CreateWorkspaceRequest = &request
	_, err := csm.workspaceCient.CreateWorkspace(&params, csm.authInfoWriter)

	if err != nil {
		csmError, ok := err.(*workspace.CreateWorkspaceDefault)
		if !ok {
			return err
		}
		return fmt.Errorf(*csmError.Payload.Message)
	}

	return nil

}
func (csm *csmClient) WorkspaceExists(workspaceID string) (bool, error) {
	if !csm.loggedIn {
		return false, errors.New("Not logged in")
	}
	csm.logger.Info("csm-workspace-exists", lager.Data{"workspaceID": workspaceID})

	params := workspace.GetWorkspaceParams{}
	params.WorkspaceID = workspaceID
	_, err := csm.workspaceCient.GetWorkspace(&params, csm.authInfoWriter)

	if err != nil {
		csmError, ok := err.(*workspace.GetWorkspaceDefault)
		if !ok {
			return false, err
		}

		if csmError.Code() == http.StatusNotFound {
			return false, nil
		}
		return false, fmt.Errorf(*csmError.Payload.Message)
	}

	return true, nil

}
func (csm *csmClient) DeleteWorkspace(workspaceID string) error {
	if !csm.loggedIn {
		return errors.New("Not logged in")
	}
	csm.logger.Info("csm-delete-workspace", lager.Data{"workspaceID": workspaceID})
	params := workspace.DeleteWorkspaceParams{}
	params.WorkspaceID = workspaceID
	_, err := csm.workspaceCient.DeleteWorkspace(&params, csm.authInfoWriter)

	if err != nil {
		csmError, ok := err.(*workspace.DeleteWorkspaceDefault)
		if !ok {
			return err
		}
		return fmt.Errorf(*csmError.Payload.Message)
	}

	//TODO: does not throw an error if the workspace does not exist
	return nil

}
func (csm *csmClient) CreateConnection(workspaceID, connectionID string) (interface{}, error) {
	if !csm.loggedIn {
		return nil, errors.New("Not logged in")
	}
	csm.logger.Info("csm-create-connection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})
	params := connection.CreateConnectionParams{}
	params.WorkspaceID = workspaceID

	request := models.ServiceManagerConnectionCreateRequest{
		ConnectionID: &connectionID,
	}

	params.ConnectionCreateRequest = &request
	response, err := csm.connectionClient.CreateConnection(&params, csm.authInfoWriter)
	if err != nil {
		csmError, ok := err.(*connection.CreateConnectionDefault)
		if !ok {
			return nil, err
		}
		return nil, fmt.Errorf(*csmError.Payload.Message)
	}

	return response.Payload.Details, nil

}
func (csm *csmClient) ConnectionExists(workspaceID, connectionID string) (bool, error) {
	if !csm.loggedIn {
		return false, errors.New("Not logged in")
	}
	csm.logger.Info("csm-connection-exists", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})
	params := connection.GetConnectionParams{
		WorkspaceID:  workspaceID,
		ConnectionID: connectionID,
	}

	_, err := csm.connectionClient.GetConnection(&params, csm.authInfoWriter)

	if err != nil {
		csmError, ok := err.(*connection.GetConnectionDefault)
		if !ok {
			return false, err
		}
		if csmError.Code() == http.StatusNotFound {
			return false, nil
		}
		return false, fmt.Errorf(*csmError.Payload.Message)
	}

	return true, nil

}
func (csm *csmClient) DeleteConnection(workspaceID, connectionID string) error {
	if !csm.loggedIn {
		return errors.New("Not logged in")
	}
	csm.logger.Info("csm-delete-connection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})
	params := connection.DeleteConnectionParams{
		WorkspaceID:  workspaceID,
		ConnectionID: connectionID,
	}

	_, err := csm.connectionClient.DeleteConnection(&params, csm.authInfoWriter)
	if err != nil {
		csmError, ok := err.(*connection.DeleteConnectionDefault)
		if !ok {
			return err
		}
		return fmt.Errorf(*csmError.Payload.Message)
	}

	//TODO: in CSM this passes all the time, it does not take into consideration if a connection exists
	return nil

}
