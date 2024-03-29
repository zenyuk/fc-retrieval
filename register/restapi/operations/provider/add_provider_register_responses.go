// Code generated by go-swagger; DO NOT EDIT.

package provider

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/ConsenSys/fc-retrieval/register/models"
)

// AddProviderRegisterOKCode is the HTTP code returned for type AddProviderRegisterOK
const AddProviderRegisterOKCode int = 200

/*AddProviderRegisterOK Provider register added

swagger:response addProviderRegisterOK
*/
type AddProviderRegisterOK struct {

	/*
	  In: Body
	*/
	Payload *models.ProviderRegister `json:"body,omitempty"`
}

// NewAddProviderRegisterOK creates AddProviderRegisterOK with default headers values
func NewAddProviderRegisterOK() *AddProviderRegisterOK {

	return &AddProviderRegisterOK{}
}

// WithPayload adds the payload to the add provider register o k response
func (o *AddProviderRegisterOK) WithPayload(payload *models.ProviderRegister) *AddProviderRegisterOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add provider register o k response
func (o *AddProviderRegisterOK) SetPayload(payload *models.ProviderRegister) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddProviderRegisterOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*AddProviderRegisterDefault Internal error

swagger:response addProviderRegisterDefault
*/
type AddProviderRegisterDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewAddProviderRegisterDefault creates AddProviderRegisterDefault with default headers values
func NewAddProviderRegisterDefault(code int) *AddProviderRegisterDefault {
	if code <= 0 {
		code = 500
	}

	return &AddProviderRegisterDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the add provider register default response
func (o *AddProviderRegisterDefault) WithStatusCode(code int) *AddProviderRegisterDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the add provider register default response
func (o *AddProviderRegisterDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the add provider register default response
func (o *AddProviderRegisterDefault) WithPayload(payload *models.Error) *AddProviderRegisterDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add provider register default response
func (o *AddProviderRegisterDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddProviderRegisterDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
