package service_instances

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/hpcloud/cf-usb/lib/brokermodel"
)

/*ServiceUnbindOK Binding was deleted. The expected response body is {}.

swagger:response serviceUnbindOK
*/
type ServiceUnbindOK struct {

	// In: body
	Payload *brokermodel.Empty `json:"body,omitempty"`
}

// NewServiceUnbindOK creates ServiceUnbindOK with default headers values
func NewServiceUnbindOK() *ServiceUnbindOK {
	return &ServiceUnbindOK{}
}

// WithPayload adds the payload to the service unbind o k response
func (o *ServiceUnbindOK) WithPayload(payload *brokermodel.Empty) *ServiceUnbindOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the service unbind o k response
func (o *ServiceUnbindOK) SetPayload(payload *brokermodel.Empty) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ServiceUnbindOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*ServiceUnbindGone Should be returned if the binding does not exist. The expected response body is {}.

swagger:response serviceUnbindGone
*/
type ServiceUnbindGone struct {

	// In: body
	Payload *brokermodel.Empty `json:"body,omitempty"`
}

// NewServiceUnbindGone creates ServiceUnbindGone with default headers values
func NewServiceUnbindGone() *ServiceUnbindGone {
	return &ServiceUnbindGone{}
}

// WithPayload adds the payload to the service unbind gone response
func (o *ServiceUnbindGone) WithPayload(payload *brokermodel.Empty) *ServiceUnbindGone {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the service unbind gone response
func (o *ServiceUnbindGone) SetPayload(payload *brokermodel.Empty) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ServiceUnbindGone) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(410)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*ServiceUnbindDefault generic error response

swagger:response serviceUnbindDefault
*/
type ServiceUnbindDefault struct {
	_statusCode int

	// In: body
	Payload *brokermodel.BrokerError `json:"body,omitempty"`
}

// NewServiceUnbindDefault creates ServiceUnbindDefault with default headers values
func NewServiceUnbindDefault(code int) *ServiceUnbindDefault {
	if code <= 0 {
		code = 500
	}

	return &ServiceUnbindDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the service unbind default response
func (o *ServiceUnbindDefault) WithStatusCode(code int) *ServiceUnbindDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the service unbind default response
func (o *ServiceUnbindDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the service unbind default response
func (o *ServiceUnbindDefault) WithPayload(payload *brokermodel.BrokerError) *ServiceUnbindDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the service unbind default response
func (o *ServiceUnbindDefault) SetPayload(payload *brokermodel.BrokerError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ServiceUnbindDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
