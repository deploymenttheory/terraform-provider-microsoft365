package sharedmodels

// ResourceIdentity provides a shared struct for resources implementing ResourceWithIdentity.
//
// The Terraform Plugin Framework separates resource identity (minimal ID for import/list operations)
// from full resource state. The framework requires identity to be set as a struct matching the
// IdentitySchema shape, not as a raw value. This shared struct ensures type safety and consistency
// across all resources when calling resp.Identity.Set() in Read operations.
type ResourceIdentity struct {
	ID string `tfsdk:"id"`
}
