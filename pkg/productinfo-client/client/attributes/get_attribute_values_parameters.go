// Code generated by go-swagger; DO NOT EDIT.

package attributes

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetAttributeValuesParams creates a new GetAttributeValuesParams object
// with the default values initialized.
func NewGetAttributeValuesParams() *GetAttributeValuesParams {
	var ()
	return &GetAttributeValuesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetAttributeValuesParamsWithTimeout creates a new GetAttributeValuesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetAttributeValuesParamsWithTimeout(timeout time.Duration) *GetAttributeValuesParams {
	var ()
	return &GetAttributeValuesParams{

		timeout: timeout,
	}
}

// NewGetAttributeValuesParamsWithContext creates a new GetAttributeValuesParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetAttributeValuesParamsWithContext(ctx context.Context) *GetAttributeValuesParams {
	var ()
	return &GetAttributeValuesParams{

		Context: ctx,
	}
}

// NewGetAttributeValuesParamsWithHTTPClient creates a new GetAttributeValuesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetAttributeValuesParamsWithHTTPClient(client *http.Client) *GetAttributeValuesParams {
	var ()
	return &GetAttributeValuesParams{
		HTTPClient: client,
	}
}

/*GetAttributeValuesParams contains all the parameters to send to the API endpoint
for the get attribute values operation typically these are written to a http.Request
*/
type GetAttributeValuesParams struct {

	/*Attribute*/
	Attribute string
	/*Provider*/
	Provider string
	/*Region*/
	Region string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get attribute values params
func (o *GetAttributeValuesParams) WithTimeout(timeout time.Duration) *GetAttributeValuesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get attribute values params
func (o *GetAttributeValuesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get attribute values params
func (o *GetAttributeValuesParams) WithContext(ctx context.Context) *GetAttributeValuesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get attribute values params
func (o *GetAttributeValuesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get attribute values params
func (o *GetAttributeValuesParams) WithHTTPClient(client *http.Client) *GetAttributeValuesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get attribute values params
func (o *GetAttributeValuesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAttribute adds the attribute to the get attribute values params
func (o *GetAttributeValuesParams) WithAttribute(attribute string) *GetAttributeValuesParams {
	o.SetAttribute(attribute)
	return o
}

// SetAttribute adds the attribute to the get attribute values params
func (o *GetAttributeValuesParams) SetAttribute(attribute string) {
	o.Attribute = attribute
}

// WithProvider adds the provider to the get attribute values params
func (o *GetAttributeValuesParams) WithProvider(provider string) *GetAttributeValuesParams {
	o.SetProvider(provider)
	return o
}

// SetProvider adds the provider to the get attribute values params
func (o *GetAttributeValuesParams) SetProvider(provider string) {
	o.Provider = provider
}

// WithRegion adds the region to the get attribute values params
func (o *GetAttributeValuesParams) WithRegion(region string) *GetAttributeValuesParams {
	o.SetRegion(region)
	return o
}

// SetRegion adds the region to the get attribute values params
func (o *GetAttributeValuesParams) SetRegion(region string) {
	o.Region = region
}

// WriteToRequest writes these params to a swagger request
func (o *GetAttributeValuesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param attribute
	if err := r.SetPathParam("attribute", o.Attribute); err != nil {
		return err
	}

	// path param provider
	if err := r.SetPathParam("provider", o.Provider); err != nil {
		return err
	}

	// path param region
	if err := r.SetPathParam("region", o.Region); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
