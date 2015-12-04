package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit"

	"github.com/hpcloud/cf-usb/lib/genmodel"
)

/*GetInfoOK Successful response

swagger:response getInfoOK
*/
type GetInfoOK struct {

	// In: body
	Payload *genmodel.Info `json:"body,omitempty"`
}

// WriteResponse to the client
func (o *GetInfoOK) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetInfoInternalServerError Unexpected error

swagger:response getInfoInternalServerError
*/
type GetInfoInternalServerError struct {

	// In: body
	Payload string `json:"body,omitempty"`
}

// WriteResponse to the client
func (o *GetInfoInternalServerError) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(500)
	if err := producer.Produce(rw, o.Payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
