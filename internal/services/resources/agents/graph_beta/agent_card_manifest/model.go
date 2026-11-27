package graphBetaAgentCardManifest

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var errInvalidManifest = errors.New("invalid agent card manifest")

type AgentCardManifestResourceModel struct {
	ID                                types.String         `tfsdk:"id"`
	AgentCardManifestID               types.String         `tfsdk:"agent_card_manifest_id"`
	DisplayName                       types.String         `tfsdk:"display_name"`
	Description                       types.String         `tfsdk:"description"`
	IconURL                           types.String         `tfsdk:"icon_url"`
	ProtocolVersion                   types.String         `tfsdk:"protocol_version"`
	Version                           types.String         `tfsdk:"version"`
	DocumentationURL                  types.String         `tfsdk:"documentation_url"`
	ManagedBy                         types.String         `tfsdk:"managed_by"`
	OriginatingStore                  types.String         `tfsdk:"originating_store"`
	OwnerIDs                          types.Set            `tfsdk:"owner_ids"`
	DefaultInputModes                 types.Set            `tfsdk:"default_input_modes"`
	DefaultOutputModes                types.Set            `tfsdk:"default_output_modes"`
	SupportsAuthenticatedExtendedCard types.Bool           `tfsdk:"supports_authenticated_extended_card"`
	Provider                          types.Object         `tfsdk:"agent_provider"`
	Capabilities                      types.Object         `tfsdk:"capabilities"`
	Skills                            types.List           `tfsdk:"skills"`
	SecuritySchemesJSON               jsontypes.Normalized `tfsdk:"security_schemes_json"`
	SecurityJSON                      jsontypes.Normalized `tfsdk:"security_json"`
	CreatedBy                         types.String         `tfsdk:"created_by"`
	CreatedDateTime                   types.String         `tfsdk:"created_date_time"`
	LastModifiedDateTime              types.String         `tfsdk:"last_modified_date_time"`
	Timeouts                          timeouts.Value       `tfsdk:"timeouts"`
}
type ProviderModel struct {
	Organization types.String `tfsdk:"organization"`
	URL          types.String `tfsdk:"url"`
}
type CapabilitiesModel struct {
	Streaming              types.Bool `tfsdk:"streaming"`
	PushNotifications      types.Bool `tfsdk:"push_notifications"`
	StateTransitionHistory types.Bool `tfsdk:"state_transition_history"`
	Extensions             types.List `tfsdk:"extensions"`
}
type ExtensionModel struct {
	URI         types.String         `tfsdk:"uri"`
	Description types.String         `tfsdk:"description"`
	Required    types.Bool           `tfsdk:"required"`
	ParamsJSON  jsontypes.Normalized `tfsdk:"params_json"`
}
type SkillModel struct {
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Tags        types.Set    `tfsdk:"tags"`
	Examples    types.Set    `tfsdk:"examples"`
	InputModes  types.Set    `tfsdk:"input_modes"`
	OutputModes types.Set    `tfsdk:"output_modes"`
}

var (
	providerAttrTypes     = map[string]attr.Type{"organization": types.StringType, "url": types.StringType}
	extensionAttrTypes    = map[string]attr.Type{"uri": types.StringType, "description": types.StringType, "required": types.BoolType, "params_json": jsontypes.NormalizedType{}}
	capabilitiesAttrTypes = map[string]attr.Type{
		"streaming":                types.BoolType,
		"push_notifications":       types.BoolType,
		"state_transition_history": types.BoolType,
		"extensions":               types.ListType{ElemType: types.ObjectType{AttrTypes: extensionAttrTypes}},
	}
	skillAttrTypes = map[string]attr.Type{
		"id":           types.StringType,
		"display_name": types.StringType,
		"description":  types.StringType,
		"tags":         types.SetType{ElemType: types.StringType},
		"examples":     types.SetType{ElemType: types.StringType},
		"input_modes":  types.SetType{ElemType: types.StringType},
		"output_modes": types.SetType{ElemType: types.StringType},
	}
)
