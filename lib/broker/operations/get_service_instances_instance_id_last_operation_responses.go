package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/SUSE/cf-usb/lib/brokermodel"
)

/*GetServiceInstancesInstanceIDLastOperationOK Valid state values are 'in progress', 'succeeded', and 'failed'

swagger:response getServiceInstancesInstanceIdLastOperationOK
*/
type GetServiceInstancesInstanceIDLastOperationOK struct {

	// In: body
	Payload *brokermodel.LastOperation `json:"body,omitempty"`
}

// NewGetServiceInstancesInstanceIDLastOperationOK creates GetServiceInstancesInstanceIDLastOperationOK with default headers values
func NewGetServiceInstancesInstanceIDLastOperationOK() *GetServiceInstancesInstanceIDLastOperationOK {
	return &GetServiceInstancesInstanceIDLastOperationOK{}
}

// WithPayload adds the payload to the get service instances instance Id last operation o k response
func (o *GetServiceInstancesInstanceIDLastOperationOK) WithPayload(payload *brokermodel.LastOperation) *GetServiceInstancesInstanceIDLastOperationOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get service instances instance Id last operation o k response
func (o *GetServiceInstancesInstanceIDLastOperationOK) SetPayload(payload *brokermodel.LastOperation) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetServiceInstancesInstanceIDLastOperationOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetServiceInstancesInstanceIDLastOperationGone Appropriate only for asynchronous delete requests. Cloud Foundry will consider this response a success and remove the resource from its database.

swagger:response getServiceInstancesInstanceIdLastOperationGone
*/
type GetServiceInstancesInstanceIDLastOperationGone struct {

	// In: body
	Payload brokermodel.Empty `json:"body,omitempty"`
}

// NewGetServiceInstancesInstanceIDLastOperationGone creates GetServiceInstancesInstanceIDLastOperationGone with default headers values
func NewGetServiceInstancesInstanceIDLastOperationGone() *GetServiceInstancesInstanceIDLastOperationGone {
	return &GetServiceInstancesInstanceIDLastOperationGone{}
}

// WithPayload adds the payload to the get service instances instance Id last operation gone response
func (o *GetServiceInstancesInstanceIDLastOperationGone) WithPayload(payload brokermodel.Empty) *GetServiceInstancesInstanceIDLastOperationGone {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get service instances instance Id last operation gone response
func (o *GetServiceInstancesInstanceIDLastOperationGone) SetPayload(payload brokermodel.Empty) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetServiceInstancesInstanceIDLastOperationGone) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(410)
	if err := producer.Produce(rw, o.Payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

/*GetServiceInstancesInstanceIDLastOperationDefault generic error response

swagger:response getServiceInstancesInstanceIdLastOperationDefault
*/
type GetServiceInstancesInstanceIDLastOperationDefault struct {
	_statusCode int

	// In: body
	Payload *brokermodel.BrokerError `json:"body,omitempty"`
}

// NewGetServiceInstancesInstanceIDLastOperationDefault creates GetServiceInstancesInstanceIDLastOperationDefault with default headers values
func NewGetServiceInstancesInstanceIDLastOperationDefault(code int) *GetServiceInstancesInstanceIDLastOperationDefault {
	if code <= 0 {
		code = 500
	}

	return &GetServiceInstancesInstanceIDLastOperationDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get service instances instance ID last operation default response
func (o *GetServiceInstancesInstanceIDLastOperationDefault) WithStatusCode(code int) *GetServiceInstancesInstanceIDLastOperationDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get service instances instance ID last operation default response
func (o *GetServiceInstancesInstanceIDLastOperationDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get service instances instance ID last operation default response
func (o *GetServiceInstancesInstanceIDLastOperationDefault) WithPayload(payload *brokermodel.BrokerError) *GetServiceInstancesInstanceIDLastOperationDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get service instances instance ID last operation default response
func (o *GetServiceInstancesInstanceIDLastOperationDefault) SetPayload(payload *brokermodel.BrokerError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetServiceInstancesInstanceIDLastOperationDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
