package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/SUSE/cf-usb/lib/brokermodel"
)

/*ServiceBindOK May be returned if the binding already exists and the requested parameters are identical to the existing binding.

swagger:response serviceBindOK
*/
type ServiceBindOK struct {

	// In: body
	Payload *brokermodel.BindingResponse `json:"body,omitempty"`
}

// NewServiceBindOK creates ServiceBindOK with default headers values
func NewServiceBindOK() *ServiceBindOK {
	return &ServiceBindOK{}
}

// WithPayload adds the payload to the service bind o k response
func (o *ServiceBindOK) WithPayload(payload *brokermodel.BindingResponse) *ServiceBindOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the service bind o k response
func (o *ServiceBindOK) SetPayload(payload *brokermodel.BindingResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ServiceBindOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*ServiceBindCreated Binding has been created.

swagger:response serviceBindCreated
*/
type ServiceBindCreated struct {

	// In: body
	Payload *brokermodel.BindingResponse `json:"body,omitempty"`
}

// NewServiceBindCreated creates ServiceBindCreated with default headers values
func NewServiceBindCreated() *ServiceBindCreated {
	return &ServiceBindCreated{}
}

// WithPayload adds the payload to the service bind created response
func (o *ServiceBindCreated) WithPayload(payload *brokermodel.BindingResponse) *ServiceBindCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the service bind created response
func (o *ServiceBindCreated) SetPayload(payload *brokermodel.BindingResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ServiceBindCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*ServiceBindConflict Should be returned if the requested binding already exists. The expected response body is {}, though the description field can be used to return a user-facing error message, as described in Broker Errors.

swagger:response serviceBindConflict
*/
type ServiceBindConflict struct {

	// In: body
	Payload brokermodel.Empty `json:"body,omitempty"`
}

// NewServiceBindConflict creates ServiceBindConflict with default headers values
func NewServiceBindConflict() *ServiceBindConflict {
	return &ServiceBindConflict{}
}

// WithPayload adds the payload to the service bind conflict response
func (o *ServiceBindConflict) WithPayload(payload brokermodel.Empty) *ServiceBindConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the service bind conflict response
func (o *ServiceBindConflict) SetPayload(payload brokermodel.Empty) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ServiceBindConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	if err := producer.Produce(rw, o.Payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

/*ServiceBindUnprocessableEntity Should be returned if the broker requires that app_guid be included in the request body.

swagger:response serviceBindUnprocessableEntity
*/
type ServiceBindUnprocessableEntity struct {

	// In: body
	Payload brokermodel.Empty `json:"body,omitempty"`
}

// NewServiceBindUnprocessableEntity creates ServiceBindUnprocessableEntity with default headers values
func NewServiceBindUnprocessableEntity() *ServiceBindUnprocessableEntity {
	return &ServiceBindUnprocessableEntity{}
}

// WithPayload adds the payload to the service bind unprocessable entity response
func (o *ServiceBindUnprocessableEntity) WithPayload(payload brokermodel.Empty) *ServiceBindUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the service bind unprocessable entity response
func (o *ServiceBindUnprocessableEntity) SetPayload(payload brokermodel.Empty) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ServiceBindUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if err := producer.Produce(rw, o.Payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

/*ServiceBindDefault generic error response

swagger:response serviceBindDefault
*/
type ServiceBindDefault struct {
	_statusCode int

	// In: body
	Payload *brokermodel.BrokerError `json:"body,omitempty"`
}

// NewServiceBindDefault creates ServiceBindDefault with default headers values
func NewServiceBindDefault(code int) *ServiceBindDefault {
	if code <= 0 {
		code = 500
	}

	return &ServiceBindDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the service bind default response
func (o *ServiceBindDefault) WithStatusCode(code int) *ServiceBindDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the service bind default response
func (o *ServiceBindDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the service bind default response
func (o *ServiceBindDefault) WithPayload(payload *brokermodel.BrokerError) *ServiceBindDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the service bind default response
func (o *ServiceBindDefault) SetPayload(payload *brokermodel.BrokerError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ServiceBindDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
