package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit"

	"github.com/hpcloud/cf-usb/lib/genmodel"
)

/*CreateDriverCreated Driver created

swagger:response createDriverCreated
*/
type CreateDriverCreated struct {

	// In: body
	Payload *genmodel.Driver `json:"body,omitempty"`
}

// WriteResponse to the client
func (o *CreateDriverCreated) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*CreateDriverConflict A driver with the same type already exists

swagger:response createDriverConflict
*/
type CreateDriverConflict struct {
}

// WriteResponse to the client
func (o *CreateDriverConflict) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(409)
}

/*CreateDriverInternalServerError Unexpected error

swagger:response createDriverInternalServerError
*/
type CreateDriverInternalServerError struct {

	// In: body
	Payload string `json:"body,omitempty"`
}

// WriteResponse to the client
func (o *CreateDriverInternalServerError) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(500)
	if err := producer.Produce(rw, o.Payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
