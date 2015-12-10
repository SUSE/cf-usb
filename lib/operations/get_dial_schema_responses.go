package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit"

	"github.com/hpcloud/cf-usb/lib/genmodel"
)

/*GetDialSchemaOK Sucessfull response

swagger:response getDialSchemaOK
*/
type GetDialSchemaOK struct {

	// In: body
	Payload genmodel.DialSchema `json:"body,omitempty"`
}

// WriteResponse to the client
func (o *GetDialSchemaOK) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(200)
	if err := producer.Produce(rw, o.Payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

/*GetDialSchemaNotFound Not Found

swagger:response getDialSchemaNotFound
*/
type GetDialSchemaNotFound struct {
}

// WriteResponse to the client
func (o *GetDialSchemaNotFound) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(404)
}

/*GetDialSchemaInternalServerError Unexpected error

swagger:response getDialSchemaInternalServerError
*/
type GetDialSchemaInternalServerError struct {

	// In: body
	Payload string `json:"body,omitempty"`
}

// WriteResponse to the client
func (o *GetDialSchemaInternalServerError) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(500)
	if err := producer.Produce(rw, o.Payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}