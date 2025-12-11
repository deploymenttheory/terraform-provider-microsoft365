// REF: https://learn.microsoft.com/en-us/graph/api/resources/agentinstance?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentregistry-post-agentinstances?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentinstance-get?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentinstance-update?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentregistry-delete-agentinstances?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentcardmanifest-get?view=graph-rest-beta
package graphBetaAgentInstance

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AgentInstanceResourceModel represents the schema for the Agent Instance resource
type AgentInstanceResourceModel struct {
	ID                       types.String                    `tfsdk:"id"`
	DisplayName              types.String                    `tfsdk:"display_name"`
	OwnerIds                 types.Set                       `tfsdk:"owner_ids"`
	ManagedBy                types.String                    `tfsdk:"managed_by"`
	OriginatingStore         types.String                    `tfsdk:"originating_store"`
	CreatedBy                types.String                    `tfsdk:"created_by"`
	SourceAgentId            types.String                    `tfsdk:"source_agent_id"`
	AgentIdentityBlueprintId types.String                    `tfsdk:"agent_identity_blueprint_id"`
	AgentIdentityId          types.String                    `tfsdk:"agent_identity_id"`
	AgentUserId              types.String                    `tfsdk:"agent_user_id"`
	CreatedDateTime          types.String                    `tfsdk:"created_date_time"`
	LastModifiedDateTime     types.String                    `tfsdk:"last_modified_date_time"`
	Url                      types.String                    `tfsdk:"url"`
	PreferredTransport       types.String                    `tfsdk:"preferred_transport"`
	AdditionalInterfaces     []AgentInterfaceModel           `tfsdk:"additional_interfaces"`
	AgentCardManifest        *AgentCardManifestResourceModel `tfsdk:"agent_card_manifest"`
	Timeouts                 timeouts.Value                  `tfsdk:"timeouts"`
}

// AgentInterfaceModel represents an additional interface/transport supported by the agent
type AgentInterfaceModel struct {
	Url       types.String `tfsdk:"url"`
	Transport types.String `tfsdk:"transport"`
}

// AgentCardManifestResourceModel represents the agent card manifest associated with an agent instance
type AgentCardManifestResourceModel struct {
	ID                                types.String                `tfsdk:"id"`
	OwnerIds                          types.Set                   `tfsdk:"owner_ids"`
	OriginatingStore                  types.String                `tfsdk:"originating_store"`
	DisplayName                       types.String                `tfsdk:"display_name"`
	Description                       types.String                `tfsdk:"description"`
	IconUrl                           types.String                `tfsdk:"icon_url"`
	Provider                          *AgentCardProviderModel     `tfsdk:"provider"`
	ProtocolVersion                   types.String                `tfsdk:"protocol_version"`
	Version                           types.String                `tfsdk:"version"`
	DocumentationUrl                  types.String                `tfsdk:"documentation_url"`
	Capabilities                      *AgentCardCapabilitiesModel `tfsdk:"capabilities"`
	DefaultInputModes                 types.Set                   `tfsdk:"default_input_modes"`
	DefaultOutputModes                types.Set                   `tfsdk:"default_output_modes"`
	SupportsAuthenticatedExtendedCard types.Bool                  `tfsdk:"supports_authenticated_extended_card"`
	Skills                            types.List                  `tfsdk:"skills"`
}

// AgentCardProviderModel represents the provider information for an agent card
type AgentCardProviderModel struct {
	Organization types.String `tfsdk:"organization"`
	Url          types.String `tfsdk:"url"`
}

// AgentCardCapabilitiesModel represents the capabilities of an agent card
type AgentCardCapabilitiesModel struct {
	Streaming              types.Bool `tfsdk:"streaming"`
	PushNotifications      types.Bool `tfsdk:"push_notifications"`
	StateTransitionHistory types.Bool `tfsdk:"state_transition_history"`
	Extensions             types.List `tfsdk:"extensions"`
}

// AgentCardCapabilityExtensionModel represents a capability extension
type AgentCardCapabilityExtensionModel struct {
	Uri         types.String `tfsdk:"uri"`
	Description types.String `tfsdk:"description"`
	Required    types.Bool   `tfsdk:"required"`
	Params      types.Map    `tfsdk:"params"`
}

// AgentCardSkillModel represents a skill defined in the agent card manifest
type AgentCardSkillModel struct {
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Tags        types.Set    `tfsdk:"tags"`
	Examples    types.Set    `tfsdk:"examples"`
	InputModes  types.Set    `tfsdk:"input_modes"`
	OutputModes types.Set    `tfsdk:"output_modes"`
}
