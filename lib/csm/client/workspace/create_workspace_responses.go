package workspace

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/hpcloud/cf-usb/lib/csm/models"
)

// CreateWorkspaceReader is a Reader for the CreateWorkspace structure.
type CreateWorkspaceReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the recieved o.
func (o *CreateWorkspaceReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 201:
		result := NewCreateWorkspaceCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewCreateWorkspaceDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	}
}

// NewCreateWorkspaceCreated creates a CreateWorkspaceCreated with default headers values
func NewCreateWorkspaceCreated() *CreateWorkspaceCreated {
	return &CreateWorkspaceCreated{}
}

/*CreateWorkspaceCreated handles this case with default header values.

create workspace
*/
type CreateWorkspaceCreated struct {
	Payload *models.ServiceManagerWorkspaceResponse
}

func (o *CreateWorkspaceCreated) Error() string {
	return fmt.Sprintf("[POST /workspaces][%d] createWorkspaceCreated  %+v", 201, o.Payload)
}

func (o *CreateWorkspaceCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ServiceManagerWorkspaceResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateWorkspaceDefault creates a CreateWorkspaceDefault with default headers values
func NewCreateWorkspaceDefault(code int) *CreateWorkspaceDefault {
	return &CreateWorkspaceDefault{
		_statusCode: code,
	}
}

/*CreateWorkspaceDefault handles this case with default header values.

generic error response
*/
type CreateWorkspaceDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the create workspace default response
func (o *CreateWorkspaceDefault) Code() int {
	return o._statusCode
}

func (o *CreateWorkspaceDefault) Error() string {
	return fmt.Sprintf("[POST /workspaces][%d] createWorkspace default  %+v", o._statusCode, o.Payload)
}

func (o *CreateWorkspaceDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}