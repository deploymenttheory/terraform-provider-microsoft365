// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementintent?view=graph-rest-beta
package graphBetaEndpointSecurityIntent

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// EndpointSecurityIntentResourceModel represents the schema for the EndpointSecurityIntent resource
type EndpointSecurityIntentResourceModel struct {
	ID                               types.String   `tfsdk:"id"`
	TemplateId                       types.String   `tfsdk:"template_id"`
	DisplayName                      types.String   `tfsdk:"display_name"`
	Description                      types.String   `tfsdk:"description"`
	RoleScopeTagIds                  types.Set      `tfsdk:"role_scope_tag_ids"`
	Settings                         types.Set      `tfsdk:"settings"`
	Assignments                      types.Set      `tfsdk:"assignments"`
	IsAssigned                       types.Bool     `tfsdk:"is_assigned"`
	IsMigratingToConfigurationPolicy types.Bool     `tfsdk:"is_migrating_to_configuration_policy"`
	LastModifiedDateTime             types.String   `tfsdk:"last_modified_date_time"`
	Timeouts                         timeouts.Value `tfsdk:"timeouts"`
}

// EndpointSecurityIntentSettingModel represents one template setting managed by the resource
type EndpointSecurityIntentSettingModel struct {
	DefinitionId types.String `tfsdk:"definition_id"`
	SettingType  types.String `tfsdk:"setting_type"`
	ValueJson    types.String `tfsdk:"value_json"`
}
