// REF: https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicy?view=graph-rest-beta
package graphBetaRoleManagementPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RoleManagementPolicyResourceModel represents the schema for the RoleManagementPolicy resource
type RoleManagementPolicyResourceModel struct {
	ID                         types.String                              `tfsdk:"id"`
	PolicyId                   types.String                              `tfsdk:"policy_id"`
	RoleDefinitionId           types.String                              `tfsdk:"role_definition_id"`
	ScopeId                    types.String                              `tfsdk:"scope_id"`
	ScopeType                  types.String                              `tfsdk:"scope_type"`
	DisplayName                types.String                              `tfsdk:"display_name"`
	Description                types.String                              `tfsdk:"description"`
	IsOrganizationDefault      types.Bool                                `tfsdk:"is_organization_default"`
	LastModifiedDateTime       types.String                              `tfsdk:"last_modified_date_time"`
	ExpirationRules            map[string]ExpirationRuleModel            `tfsdk:"expiration_rules"`
	EnablementRules            map[string]EnablementRuleModel            `tfsdk:"enablement_rules"`
	NotificationRules          map[string]NotificationRuleModel          `tfsdk:"notification_rules"`
	AuthenticationContextRules map[string]AuthenticationContextRuleModel `tfsdk:"authentication_context_rules"`
	Timeouts                   timeouts.Value                            `tfsdk:"timeouts"`
}

// ExpirationRuleModel represents a unifiedRoleManagementPolicyExpirationRule
type ExpirationRuleModel struct {
	IsExpirationRequired types.Bool   `tfsdk:"is_expiration_required"`
	MaximumDuration      types.String `tfsdk:"maximum_duration"`
}

// EnablementRuleModel represents a unifiedRoleManagementPolicyEnablementRule
type EnablementRuleModel struct {
	EnabledRules types.Set `tfsdk:"enabled_rules"`
}

// NotificationRuleModel represents a unifiedRoleManagementPolicyNotificationRule
type NotificationRuleModel struct {
	NotificationLevel          types.String `tfsdk:"notification_level"`
	IsDefaultRecipientsEnabled types.Bool   `tfsdk:"is_default_recipients_enabled"`
	NotificationRecipients     types.Set    `tfsdk:"notification_recipients"`
}

// AuthenticationContextRuleModel represents a unifiedRoleManagementPolicyAuthenticationContextRule
type AuthenticationContextRuleModel struct {
	IsEnabled  types.Bool   `tfsdk:"is_enabled"`
	ClaimValue types.String `tfsdk:"claim_value"`
}
