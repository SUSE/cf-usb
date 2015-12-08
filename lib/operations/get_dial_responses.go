package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit"

	"github.com/hpcloud/cf-usb/lib/genmodel"
)

/*GetDialOK Sucessfull response

swagger:response getDialOK
*/
type GetDialOK struct {

	// In: body
	Payload *genmodel.Dial `json:"body,omitempty"`
}

// WriteResponse to the client
func (o *GetDialOK) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetDialNotFound Not Found

swagger:response getDialNotFound
*/
type GetDialNotFound struct {
}

// WriteResponse to the client
func (o *GetDialNotFound) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(404)
}

/*GetDialInternalServerError Unexpected error

swagger:response getDialInternalServerError
*/
type GetDialInternalServerError struct {

	// In: body
	Payload string `json:"body,omitempty"`
}

// WriteResponse to the client
func (o *GetDialInternalServerError) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(500)
	if err := producer.Produce(rw, o.Payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
