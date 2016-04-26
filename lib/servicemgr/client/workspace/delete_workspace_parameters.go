package workspace

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// NewDeleteWorkspaceParams creates a new DeleteWorkspaceParams object
// with the default values initialized.
func NewDeleteWorkspaceParams() *DeleteWorkspaceParams {
	var ()
	return &DeleteWorkspaceParams{}
}

/*DeleteWorkspaceParams contains all the parameters to send to the API endpoint
for the delete workspace operation typically these are written to a http.Request
*/
type DeleteWorkspaceParams struct {

	/*WorkspaceID
	  Workspace ID

	*/
	WorkspaceID string
}

// WithWorkspaceID adds the workspaceId to the delete workspace params
func (o *DeleteWorkspaceParams) WithWorkspaceID(workspaceId string) *DeleteWorkspaceParams {
	o.WorkspaceID = workspaceId
	return o
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteWorkspaceParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	var res []error

	// path param workspace_id
	if err := r.SetPathParam("workspace_id", o.WorkspaceID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}