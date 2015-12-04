package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit"

	"github.com/hpcloud/cf-usb/lib/genmodel"
)

/*GetDriverInstanceOK OK

swagger:response getDriverInstanceOK
*/
type GetDriverInstanceOK struct {

	// In: body
	Payload *genmodel.DriverInstance `json:"body,omitempty"`
}

// WriteResponse to the client
func (o *GetDriverInstanceOK) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetDriverInstanceNotFound Not Found

swagger:response getDriverInstanceNotFound
*/
type GetDriverInstanceNotFound struct {
}

// WriteResponse to the client
func (o *GetDriverInstanceNotFound) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(404)
}

/*GetDriverInstanceInternalServerError Unexpected error

swagger:response getDriverInstanceInternalServerError
*/
type GetDriverInstanceInternalServerError struct {

	// In: body
	Payload string `json:"body,omitempty"`
}

// WriteResponse to the client
func (o *GetDriverInstanceInternalServerError) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(500)
	if err := producer.Produce(rw, o.Payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
