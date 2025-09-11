// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-grouppolicy-grouppolicyconfiguration?view=graph-rest-beta
package graphBetaGroupPolicyConfiguration

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GroupPolicyConfigurationResourceModel struct {
	ID                               types.String   `tfsdk:"id"`
	DisplayName                      types.String   `tfsdk:"display_name"`
	Description                      types.String   `tfsdk:"description"`
	RoleScopeTagIds                  types.Set      `tfsdk:"role_scope_tag_ids"`
	PolicyConfigurationIngestionType types.String   `tfsdk:"policy_configuration_ingestion_type"`
	CreatedDateTime                  types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime             types.String   `tfsdk:"last_modified_date_time"`
	Assignments                      types.Set      `tfsdk:"assignments"`
	Timeouts                         timeouts.Value `tfsdk:"timeouts"`
}
