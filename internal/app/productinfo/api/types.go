package api

import "github.com/banzaicloud/telescopes/pkg/productinfo"

// GetProductDetailsParams is a placeholder for the get products route's path parameters
// swagger:parameters getProductDetails
type GetProductDetailsParams struct {
	// in:path
	Provider string `json:"provider"`
	// in:path
	Region string `json:"region"`
}

// ProductDetailsResponse Api object to be mapped to product info response
// swagger:model ProductDetailsResponse
type ProductDetailsResponse struct {
	// Products represents a slice of products for a given provider (VMs with attributes and process)
	Products []productinfo.ProductDetails `json:"products"`
}

// GetRegionsParams is a placeholder for the get regions route's path parameters
// swagger:parameters getRegions
type GetRegionsParams struct {
	// in:path
	Provider string `json:"provider"`
}

// RegionResp holds the list of available regions of a cloud provider
// swagger:response RegionsResponse
type RegionResp struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// GetAttributeValuesParams is a placeholder for the get attribute values route's path parameters
// swagger:parameters getAttributeValues
type GetAttributeValuesParams struct {
	// in:path
	Provider string `json:"provider"`
	// in:path
	Region string `json:"region"`
	// in:path
	Attribute string `json:"attribute"`
}

// AttributeResponse holds attribute values
// swagger:model AttributeResponse
type AttributeResponse struct {
	AttributeName   string    `json:"attributeName"`
	AttributeValues []float64 `json:"attributeValues"`
}
