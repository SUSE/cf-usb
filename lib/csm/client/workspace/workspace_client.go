package workspace

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new workspace API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for workspace API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
CreateWorkspace Create new workspace
*/
func (a *Client) CreateWorkspace(params *CreateWorkspaceParams, authInfo runtime.ClientAuthInfoWriter) (*CreateWorkspaceCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewCreateWorkspaceParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "createWorkspace",
		Method:             "POST",
		PathPattern:        "/workspaces",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &CreateWorkspaceReader{formats: a.formats},
		AuthInfo:           authInfo,
	})
	if err != nil {
		return nil, err
	}
	return result.(*CreateWorkspaceCreated), nil
}

/*
DeleteWorkspace Delete specified workspace
*/
func (a *Client) DeleteWorkspace(params *DeleteWorkspaceParams, authInfo runtime.ClientAuthInfoWriter) (*DeleteWorkspaceOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteWorkspaceParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "deleteWorkspace",
		Method:             "DELETE",
		PathPattern:        "/workspaces/{workspace_id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DeleteWorkspaceReader{formats: a.formats},
		AuthInfo:           authInfo,
	})
	if err != nil {
		return nil, err
	}
	return result.(*DeleteWorkspaceOK), nil
}

/*
GetWorkspace Get the details for the specified
*/
func (a *Client) GetWorkspace(params *GetWorkspaceParams, authInfo runtime.ClientAuthInfoWriter) (*GetWorkspaceOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetWorkspaceParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getWorkspace",
		Method:             "GET",
		PathPattern:        "/workspaces/{workspace_id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetWorkspaceReader{formats: a.formats},
		AuthInfo:           authInfo,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetWorkspaceOK), nil
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
