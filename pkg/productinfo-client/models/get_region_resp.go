// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// GetRegionResp GetRegionResp holds the detailed description of a specific region of a cloud provider
// swagger:model GetRegionResp
type GetRegionResp struct {

	// Id
	ID string `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// zones
	Zones []string `json:"zones"`
}

// Validate validates this get region resp
func (m *GetRegionResp) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *GetRegionResp) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GetRegionResp) UnmarshalBinary(b []byte) error {
	var res GetRegionResp
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
