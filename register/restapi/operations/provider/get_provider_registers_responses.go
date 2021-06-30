// Code generated by go-swagger; DO NOT EDIT.

package provider

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/ConsenSys/fc-retrieval-register/models"
)

// GetProviderRegistersOKCode is the HTTP code returned for type GetProviderRegistersOK
const GetProviderRegistersOKCode int = 200

/*GetProviderRegistersOK Provider register list

swagger:response getProviderRegistersOK
*/
type GetProviderRegistersOK struct {

	/*
	  In: Body
	*/
	Payload []*models.ProviderRegister `json:"body,omitempty"`
}

// NewGetProviderRegistersOK creates GetProviderRegistersOK with default headers values
func NewGetProviderRegistersOK() *GetProviderRegistersOK {

	return &GetProviderRegistersOK{}
}

// WithPayload adds the payload to the get provider registers o k response
func (o *GetProviderRegistersOK) WithPayload(payload []*models.ProviderRegister) *GetProviderRegistersOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get provider registers o k response
func (o *GetProviderRegistersOK) SetPayload(payload []*models.ProviderRegister) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProviderRegistersOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*models.ProviderRegister, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*GetProviderRegistersDefault Internal error

swagger:response getProviderRegistersDefault
*/
type GetProviderRegistersDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetProviderRegistersDefault creates GetProviderRegistersDefault with default headers values
func NewGetProviderRegistersDefault(code int) *GetProviderRegistersDefault {
	if code <= 0 {
		code = 500
	}

	return &GetProviderRegistersDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get provider registers default response
func (o *GetProviderRegistersDefault) WithStatusCode(code int) *GetProviderRegistersDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get provider registers default response
func (o *GetProviderRegistersDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get provider registers default response
func (o *GetProviderRegistersDefault) WithPayload(payload *models.Error) *GetProviderRegistersDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get provider registers default response
func (o *GetProviderRegistersDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProviderRegistersDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
