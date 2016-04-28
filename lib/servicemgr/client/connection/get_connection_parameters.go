package connection

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetConnectionParams creates a new GetConnectionParams object
// with the default values initialized.
func NewGetConnectionParams() *GetConnectionParams {
	var ()
	return &GetConnectionParams{}
}

/*GetConnectionParams contains all the parameters to send to the API endpoint
for the get connection operation typically these are written to a http.Request
*/
type GetConnectionParams struct {

	/*ConnectionID
	  connection ID

	*/
	ConnectionID string
	/*WorkspaceID
	  Workspace ID

	*/
	WorkspaceID string
}

// WithConnectionID adds the connectionId to the get connection params
func (o *GetConnectionParams) WithConnectionID(ConnectionID string) *GetConnectionParams {
	o.ConnectionID = ConnectionID
	return o
}

// WithWorkspaceID adds the workspaceId to the get connection params
func (o *GetConnectionParams) WithWorkspaceID(WorkspaceID string) *GetConnectionParams {
	o.WorkspaceID = WorkspaceID
	return o
}

// WriteToRequest writes these params to a swagger request
func (o *GetConnectionParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	var res []error

	// path param connection_id
	if err := r.SetPathParam("connection_id", o.ConnectionID); err != nil {
		return err
	}

	// path param workspace_id
	if err := r.SetPathParam("workspace_id", o.WorkspaceID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
