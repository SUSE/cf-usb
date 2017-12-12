package workspace

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/SUSE/cf-usb/lib/csm/models"
)

// NewCreateWorkspaceParams creates a new CreateWorkspaceParams object
// with the default values initialized.
func NewCreateWorkspaceParams() *CreateWorkspaceParams {
	var ()
	return &CreateWorkspaceParams{}
}

/*CreateWorkspaceParams contains all the parameters to send to the API endpoint
for the create workspace operation typically these are written to a http.Request
*/
type CreateWorkspaceParams struct {

	/*CreateWorkspaceRequest
	  The service JSON you want to post

	*/
	CreateWorkspaceRequest *models.ServiceManagerWorkspaceCreateRequest
}

// WithCreateWorkspaceRequest adds the createWorkspaceRequest to the create workspace params
func (o *CreateWorkspaceParams) WithCreateWorkspaceRequest(CreateWorkspaceRequest *models.ServiceManagerWorkspaceCreateRequest) *CreateWorkspaceParams {
	o.CreateWorkspaceRequest = CreateWorkspaceRequest
	return o
}

// WriteToRequest writes these params to a swagger request
func (o *CreateWorkspaceParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	var res []error

	if o.CreateWorkspaceRequest == nil {
		o.CreateWorkspaceRequest = new(models.ServiceManagerWorkspaceCreateRequest)
	}

	if err := r.SetBodyParam(o.CreateWorkspaceRequest); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
