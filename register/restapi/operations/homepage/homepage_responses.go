// Code generated by go-swagger; DO NOT EDIT.

package homepage

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/ConsenSys/fc-retrieval/register/models"
)

// HomepageOKCode is the HTTP code returned for type HomepageOK
const HomepageOKCode int = 200

/*HomepageOK homepage response

swagger:response homepageOK
*/
type HomepageOK struct {

	/*
	  In: Body
	*/
	Payload *models.Ack `json:"body,omitempty"`
}

// NewHomepageOK creates HomepageOK with default headers values
func NewHomepageOK() *HomepageOK {

	return &HomepageOK{}
}

// WithPayload adds the payload to the homepage o k response
func (o *HomepageOK) WithPayload(payload *models.Ack) *HomepageOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the homepage o k response
func (o *HomepageOK) SetPayload(payload *models.Ack) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *HomepageOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*HomepageDefault user validation error

swagger:response homepageDefault
*/
type HomepageDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewHomepageDefault creates HomepageDefault with default headers values
func NewHomepageDefault(code int) *HomepageDefault {
	if code <= 0 {
		code = 500
	}

	return &HomepageDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the homepage default response
func (o *HomepageDefault) WithStatusCode(code int) *HomepageDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the homepage default response
func (o *HomepageDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the homepage default response
func (o *HomepageDefault) WithPayload(payload *models.Error) *HomepageDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the homepage default response
func (o *HomepageDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *HomepageDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
