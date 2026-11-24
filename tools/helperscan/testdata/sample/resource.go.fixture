// Package sample is a fixture resource package for the helperscan tests.
package sample

import (
	"context"

	"example.com/framework/resource"
)

// SampleResource is a fixture resource.
type SampleResource struct {
	client string
}

// NewSampleResource is the standard constructor and is not reported.
func NewSampleResource() resource.Resource {
	return &SampleResource{}
}

// Metadata is a standard method and is not reported.
func (r *SampleResource) Metadata(ctx context.Context) {}

// Schema is a standard method and is not reported.
func (r *SampleResource) Schema(ctx context.Context) {}

// ImportState calls findWidgetByName directly.
func (r *SampleResource) ImportState(ctx context.Context, id string) {
	_ = findWidgetByName(nil, id)
}
