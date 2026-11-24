package sample

import "context"

// OtherResource is a second resource in the same package.
type OtherResource struct{}

// Read is a separate CRUD method from SampleResource.Read.
func (r *OtherResource) Read(ctx context.Context) {
	_ = constructWidgetRequest("other")
}
