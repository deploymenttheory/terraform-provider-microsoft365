package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Registry defines the constructors available to a provider instance.
// Test packages supply only their own constructors and explicit fixture dependencies.
type Registry struct {
	Resources          []func() resource.Resource
	DataSources        []func() datasource.DataSource
	Actions            []func() action.Action
	ListResources      []func() list.ListResource
	EphemeralResources []func() ephemeral.EphemeralResource
}

var (
	_ provider.ProviderWithEphemeralResources = &M365Provider{}
	_ provider.ProviderWithListResources      = &M365Provider{}
)

// NewWithRegistry creates a provider using the shared schema and client configuration.
// Each invocation returns an independent instance with the supplied registrations.
func NewWithRegistry(version string, unitTestMode bool, registry Registry) *M365Provider {
	p := NewMicrosoft365Provider(version, unitTestMode)().(*M365Provider)
	p.registry = registry
	return p
}

// Resources returns the resource constructors registered for this provider instance.
func (p *M365Provider) Resources(_ context.Context) []func() resource.Resource {
	return p.registry.Resources
}

// DataSources returns the data source constructors registered for this provider instance.
func (p *M365Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	return p.registry.DataSources
}

// Actions returns the action constructors registered for this provider instance.
func (p *M365Provider) Actions(_ context.Context) []func() action.Action {
	return p.registry.Actions
}

// ListResources returns the list resource constructors registered for this provider instance.
func (p *M365Provider) ListResources(_ context.Context) []func() list.ListResource {
	return p.registry.ListResources
}

// EphemeralResources returns the ephemeral resource constructors registered for this provider instance.
func (p *M365Provider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return p.registry.EphemeralResources
}
