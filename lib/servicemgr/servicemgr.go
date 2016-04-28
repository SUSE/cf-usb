package servicemgr

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	client "github.com/hpcloud/cf-usb/lib/servicemgr/client"
	connection "github.com/hpcloud/cf-usb/lib/servicemgr/client/connection"
	workspace "github.com/hpcloud/cf-usb/lib/servicemgr/client/workspace"

	models "github.com/hpcloud/cf-usb/lib/servicemgr/models"

	"github.com/pivotal-golang/lager"
)

type ServiceManager struct {
	manager *client.Servicemgr
	logger  lager.Logger
}

func NewServiceManager(transport runtime.ClientTransport, format strfmt.Registry, logger lager.Logger) ServiceManagerInterface {
	mgr := ServiceManager{}
	mgr.manager = client.New(transport, format)
	mgr.logger = logger
	return &mgr
}

func (s *ServiceManager) CreateWorkspace(request models.ServiceManagerWorkspaceCreateRequest) (models.ServiceManagerWorkspaceResponse, models.Error) {
	response := models.ServiceManagerWorkspaceResponse{}
	error := models.Error{}
	params := workspace.NewCreateWorkspaceParams()
	params.CreateWorkspaceRequest = &request

	created, err := s.manager.Workspace.CreateWorkspace(params)
	if created != nil {
		response = *created.Payload
	}
	if err != nil {
		message := err.Error()
		error.Message = &message
	}
	return response, error
}

func (s *ServiceManager) GetWorkspace(workspace_id string) (models.ServiceManagerWorkspaceResponse, models.Error) {
	s.logger.Info("get workspace", lager.Data{"id:": workspace_id})
	response := models.ServiceManagerWorkspaceResponse{}
	errorObj := models.Error{}
	params := workspace.NewGetWorkspaceParams()
	params.WorkspaceID = workspace_id

	getwork, err := s.manager.Workspace.GetWorkspace(params)
	if getwork != nil {
		if getwork.Payload != nil {
			s.logger.Info("get workspace", lager.Data{"getwork": getwork})
			response = *getwork.Payload
		}
	}
	if err != nil {
		message := err.Error()
		errorObj.Message = &message
	}
	return response, errorObj
}

func (s *ServiceManager) DeleteWorkspace(workspace_id string) models.Error {
	s.logger.Info("delete workspace", lager.Data{"workspace_id": workspace_id})

	params := workspace.NewDeleteWorkspaceParams()
	params.WorkspaceID = workspace_id
	errorObj := models.Error{}

	_, err := s.manager.Workspace.DeleteWorkspace(params)
	if err != nil {
		message := err.Error()
		errorObj.Message = &message
	}
	return errorObj
}

func (s *ServiceManager) CreateWorkspaceConnection(workspace_id string, request models.ServiceManagerConnectionCreateRequest) (models.ServiceManagerConnectionResponse, models.Error) {
	s.logger.Info("create connection", lager.Data{"workspace_id": workspace_id})
	response := models.ServiceManagerConnectionResponse{}
	error := models.Error{}
	params := connection.CreateConnectionParams{}
	params.ConnectionCreateRequest = &request
	params.WorkspaceID = workspace_id
	created, err := s.manager.Connection.CreateConnection(&params)

	if created.Payload != nil {
		response = *created.Payload
	}
	if err != nil {
		message := err.Error()
		error.Message = &message
	}
	return response, error
}

func (s *ServiceManager) GetWorkspaceConnection(workspace_id string, connection_id string) (models.ServiceManagerConnectionResponse, models.Error) {
	s.logger.Info("get connection", lager.Data{"workspace_id": workspace_id, "connection_id": connection_id})
	response := models.ServiceManagerConnectionResponse{}
	error := models.Error{}

	params := connection.GetConnectionParams{}
	params.ConnectionID = connection_id
	params.WorkspaceID = workspace_id
	created, err := s.manager.Connection.GetConnection(&params)

	if created.Payload != nil {
		response = *created.Payload
	}
	if err != nil {
		message := err.Error()
		error.Message = &message
	}

	return response, error

}

func (s *ServiceManager) DeleteWorkspaceConnection(workspace_id string, connection_id string) models.Error {
	s.logger.Info("delete connection", lager.Data{"workspace_id": workspace_id, "connection_id": connection_id})
	error := models.Error{}
	params := connection.NewDeleteConnectionParams()
	params.ConnectionID = connection_id
	params.WorkspaceID = workspace_id
	response, err := s.manager.Connection.DeleteConnection(params)
	s.logger.Debug("delete connection", lager.Data{"delete response": response, "error": err})
	if err != nil {
		message := err.Error()
		error.Message = &message
	}

	return error
}
