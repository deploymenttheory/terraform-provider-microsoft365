// REF: https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-lifecyclemanagementsettings?view=graph-rest-beta
package graphBetaIdentityAndAccessLifecycleWorkflowSettings

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// LifecycleWorkflowSettingsResourceModel represents the schema for the LifecycleWorkflowSettings resource.
// EmailSettings is a types.Object because it is Optional+Computed: when it is not configured it is unknown
// in the plan and is filled from the current tenant settings.
type LifecycleWorkflowSettingsResourceModel struct {
	ID                              types.String   `tfsdk:"id"`
	WorkflowScheduleIntervalInHours types.Int32    `tfsdk:"workflow_schedule_interval_in_hours"`
	EmailSettings                   types.Object   `tfsdk:"email_settings"`
	Timeouts                        timeouts.Value `tfsdk:"timeouts"`
}

// EmailSettingsModel represents the emailSettings used for lifecycle workflow email notifications.
type EmailSettingsModel struct {
	SenderDomain       types.String `tfsdk:"sender_domain"`
	UseCompanyBranding types.Bool   `tfsdk:"use_company_branding"`
}

var emailSettingsAttrTypes = map[string]attr.Type{
	"sender_domain":        types.StringType,
	"use_company_branding": types.BoolType,
}
