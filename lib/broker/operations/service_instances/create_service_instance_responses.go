package service_instances

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/hpcloud/cf-usb/lib/brokermodel"
)

/*CreateServiceInstanceOK May be returned if the service instance already exists and the requested parameters are identical to the existing service instance. The expected response body is below.

swagger:response createServiceInstanceOK
*/
type CreateServiceInstanceOK struct {

	// In: body
	Payload *brokermodel.DashboardURL `json:"body,omitempty"`
}

// NewCreateServiceInstanceOK creates CreateServiceInstanceOK with default headers values
func NewCreateServiceInstanceOK() *CreateServiceInstanceOK {
	return &CreateServiceInstanceOK{}
}

// WithPayload adds the payload to the create service instance o k response
func (o *CreateServiceInstanceOK) WithPayload(payload *brokermodel.DashboardURL) *CreateServiceInstanceOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create service instance o k response
func (o *CreateServiceInstanceOK) SetPayload(payload *brokermodel.DashboardURL) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateServiceInstanceOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*CreateServiceInstanceCreated Service instance has been created.

swagger:response createServiceInstanceCreated
*/
type CreateServiceInstanceCreated struct {

	// In: body
	Payload *brokermodel.DashboardURL `json:"body,omitempty"`
}

// NewCreateServiceInstanceCreated creates CreateServiceInstanceCreated with default headers values
func NewCreateServiceInstanceCreated() *CreateServiceInstanceCreated {
	return &CreateServiceInstanceCreated{}
}

// WithPayload adds the payload to the create service instance created response
func (o *CreateServiceInstanceCreated) WithPayload(payload *brokermodel.DashboardURL) *CreateServiceInstanceCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create service instance created response
func (o *CreateServiceInstanceCreated) SetPayload(payload *brokermodel.DashboardURL) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateServiceInstanceCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*CreateServiceInstanceAccepted Service instance creation has been accepted

swagger:response createServiceInstanceAccepted
*/
type CreateServiceInstanceAccepted struct {

	// In: body
	Payload *brokermodel.DashboardURL `json:"body,omitempty"`
}

// NewCreateServiceInstanceAccepted creates CreateServiceInstanceAccepted with default headers values
func NewCreateServiceInstanceAccepted() *CreateServiceInstanceAccepted {
	return &CreateServiceInstanceAccepted{}
}

// WithPayload adds the payload to the create service instance accepted response
func (o *CreateServiceInstanceAccepted) WithPayload(payload *brokermodel.DashboardURL) *CreateServiceInstanceAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create service instance accepted response
func (o *CreateServiceInstanceAccepted) SetPayload(payload *brokermodel.DashboardURL) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateServiceInstanceAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(202)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*CreateServiceInstanceConflict Should be returned if the requested service instance already exists. The expected response body is {}

swagger:response createServiceInstanceConflict
*/
type CreateServiceInstanceConflict struct {

	// In: body
	Payload brokermodel.Empty `json:"body,omitempty"`
}

// NewCreateServiceInstanceConflict creates CreateServiceInstanceConflict with default headers values
func NewCreateServiceInstanceConflict() *CreateServiceInstanceConflict {
	return &CreateServiceInstanceConflict{}
}

// WithPayload adds the payload to the create service instance conflict response
func (o *CreateServiceInstanceConflict) WithPayload(payload brokermodel.Empty) *CreateServiceInstanceConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create service instance conflict response
func (o *CreateServiceInstanceConflict) SetPayload(payload brokermodel.Empty) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateServiceInstanceConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	if err := producer.Produce(rw, o.Payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

/*CreateServiceInstanceUnprocessableEntity Shoud be returned if the broker only supports asynchronous provisioning for the requested plan and the request did not include ?accepts_incomplete=true

swagger:response createServiceInstanceUnprocessableEntity
*/
type CreateServiceInstanceUnprocessableEntity struct {

	// In: body
	Payload *brokermodel.AsyncError `json:"body,omitempty"`
}

// NewCreateServiceInstanceUnprocessableEntity creates CreateServiceInstanceUnprocessableEntity with default headers values
func NewCreateServiceInstanceUnprocessableEntity() *CreateServiceInstanceUnprocessableEntity {
	return &CreateServiceInstanceUnprocessableEntity{}
}

// WithPayload adds the payload to the create service instance unprocessable entity response
func (o *CreateServiceInstanceUnprocessableEntity) WithPayload(payload *brokermodel.AsyncError) *CreateServiceInstanceUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create service instance unprocessable entity response
func (o *CreateServiceInstanceUnprocessableEntity) SetPayload(payload *brokermodel.AsyncError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateServiceInstanceUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
