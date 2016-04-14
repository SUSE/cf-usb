package servicemgr

import (
	models "github.com/hpcloud/cf-usb/lib/servicemgr/models"
)

type ServiceManagerInterface interface {
	CreateWorkspace(models.ServiceManagerWorkspaceCreateRequest) (models.ServiceManagerWorkspaceResponse, models.Error)
	GetWorkspace(workspace_id string) (models.ServiceManagerWorkspaceResponse, models.Error)
	DeleteWorkspace(workspace_id string) models.Error
	CreateWorkspaceConnection(workspace_id string, request models.ServiceManagerConnectionCreateRequest) (models.ServiceManagerConnectionResponse, models.Error)
	GetWorkspaceConnection(workspace_id string, connection_id string) (models.ServiceManagerConnectionResponse, models.Error)
	DeleteWorkspaceConnection(workspace_id string, connection_id string) models.Error
}
