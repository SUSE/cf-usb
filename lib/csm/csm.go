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
	"github.com/hpcloud/cf-usb/lib/csm/client/status"
	"github.com/hpcloud/cf-usb/lib/csm/client/workspace"
	"github.com/hpcloud/cf-usb/lib/csm/models"
	"github.com/pivotal-golang/lager"
)

type csmClient struct {
	logger           lager.Logger
	workspaceCient   *workspace.Client
	connectionClient *connection.Client
	statusClient     *status.Client
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
	csm.statusClient = status.New(transport, strfmt.Default)
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
func (csm *csmClient) WorkspaceExists(workspaceID string) (bool, bool, error) {
	if !csm.loggedIn {
		return false, false, errors.New("Not logged in")
	}
	csm.logger.Info("csm-workspace-exists", lager.Data{"workspaceID": workspaceID})

	params := workspace.GetWorkspaceParams{}
	params.WorkspaceID = workspaceID
	response, err := csm.workspaceCient.GetWorkspace(&params, csm.authInfoWriter)

	if err != nil {
		csmError, ok := err.(*workspace.GetWorkspaceDefault)
		if !ok {
			return false, false, err
		}

		if csmError.Code() == http.StatusNotFound {
			return false, false, nil
		}
		return false, false, fmt.Errorf(*csmError.Payload.Message)
	}

	if response != nil {
		if response.Payload != nil {
			if response.Payload.ProcessingType != nil && response.Payload.Status != nil {
				if *response.Payload.ProcessingType == "none" && *response.Payload.Status == "none" {
					return false, true, nil
				}
			}
		}
	}

	return true, false, nil

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

func (csm *csmClient) ConnectionExists(workspaceID, connectionID string) (bool, bool, error) {
	if !csm.loggedIn {
		return false, false, errors.New("Not logged in")
	}
	csm.logger.Info("csm-connection-exists", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})
	params := connection.GetConnectionParams{
		WorkspaceID:  workspaceID,
		ConnectionID: connectionID,
	}

	response, err := csm.connectionClient.GetConnection(&params, csm.authInfoWriter)

	if err != nil {
		csmError, ok := err.(*connection.GetConnectionDefault)
		if !ok {
			return false, false, err
		}
		if csmError.Code() == http.StatusNotFound {
			return false, false, nil
		}
		return false, false, fmt.Errorf(*csmError.Payload.Message)
	}

	if response != nil {
		if response.Payload != nil {
			if response.Payload.ProcessingType != nil && response.Payload.Status != nil {
				if *response.Payload.ProcessingType == "none" && *response.Payload.Status == "none" {
					return false, true, nil
				}
			}
		}
	}

	return true, false, nil

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

func (csm *csmClient) GetStatus() error {
	if !csm.loggedIn {
		return errors.New("Not logged in")
	}
	params := status.NewStatusParams()
	response, err := csm.statusClient.Status(params, csm.authInfoWriter)
	if err != nil {
		csmError, ok := err.(*status.StatusDefault)
		if !ok {
			return err
		}
		return fmt.Errorf(*csmError.Payload.Message)
	}

	csm.logger.Info("status-response", lager.Data{"Status ": response.Payload.Status, "Message": response.Payload.Message, "Processing Type": response.Payload.ProcessingType})

	if response != nil {
		if *response.Payload.Status == "failed" {

			errTrace := *response.Payload.Message
			for _, diag := range response.Payload.Diagnostics {
				csm.logger.Debug("status-response-diagnostics", lager.Data{"Status": diag.Status, "Message": diag.Message, "Name": diag.Name, "Description": diag.Description})
				errTrace = errTrace + fmt.Sprintf("\n Status: %s, Name: %s, Description: %s, Message: %s", *diag.Status, *diag.Name, *diag.Description, *diag.Message)
			}

			return fmt.Errorf(errTrace)
		}
	}
	return nil
}
