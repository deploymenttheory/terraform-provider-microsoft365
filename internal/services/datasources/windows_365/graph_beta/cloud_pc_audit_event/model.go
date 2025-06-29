// REF: https://learn.microsoft.com/en-us/graph/api/resources/cloudpcauditevent?view=graph-rest-beta

package graphBetaCloudPcAuditEvent

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CloudPcAuditEventDataSourceModel represents the Terraform data source model for audit events
type CloudPcAuditEventDataSourceModel struct {
	FilterType  types.String                 `tfsdk:"filter_type"`
	FilterValue types.String                 `tfsdk:"filter_value"`
	Items       []CloudPcAuditEventItemModel `tfsdk:"items"`
	Timeouts    timeouts.Value               `tfsdk:"timeouts"`
}

type CloudPcAuditEventItemModel struct {
	ID                    types.String                `tfsdk:"id"`
	DisplayName           types.String                `tfsdk:"display_name"`
	ComponentName         types.String                `tfsdk:"component_name"`
	Activity              types.String                `tfsdk:"activity"`
	ActivityDateTime      types.String                `tfsdk:"activity_date_time"`
	ActivityType          types.String                `tfsdk:"activity_type"`
	ActivityOperationType types.String                `tfsdk:"activity_operation_type"`
	ActivityResult        types.String                `tfsdk:"activity_result"`
	CorrelationId         types.String                `tfsdk:"correlation_id"`
	Category              types.String                `tfsdk:"category"`
	Actor                 *CloudPcAuditActorModel     `tfsdk:"actor"`
	Resources             []CloudPcAuditResourceModel `tfsdk:"resources"`
}

type CloudPcAuditActorModel struct {
	ApplicationDisplayName types.String                       `tfsdk:"application_display_name"`
	ApplicationId          types.String                       `tfsdk:"application_id"`
	IpAddress              types.String                       `tfsdk:"ip_address"`
	RemoteTenantId         types.String                       `tfsdk:"remote_tenant_id"`
	RemoteUserId           types.String                       `tfsdk:"remote_user_id"`
	ServicePrincipalName   types.String                       `tfsdk:"service_principal_name"`
	Type                   types.String                       `tfsdk:"type"`
	UserId                 types.String                       `tfsdk:"user_id"`
	UserPermissions        []types.String                     `tfsdk:"user_permissions"`
	UserPrincipalName      types.String                       `tfsdk:"user_principal_name"`
	UserRoleScopeTags      []CloudPcUserRoleScopeTagInfoModel `tfsdk:"user_role_scope_tags"`
}

type CloudPcUserRoleScopeTagInfoModel struct {
	DisplayName    types.String `tfsdk:"display_name"`
	RoleScopeTagId types.String `tfsdk:"role_scope_tag_id"`
}

type CloudPcAuditResourceModel struct {
	DisplayName        types.String                `tfsdk:"display_name"`
	ModifiedProperties []CloudPcAuditPropertyModel `tfsdk:"modified_properties"`
	ResourceId         types.String                `tfsdk:"resource_id"`
	ResourceType       types.String                `tfsdk:"resource_type"`
}

type CloudPcAuditPropertyModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	NewValue    types.String `tfsdk:"new_value"`
	OldValue    types.String `tfsdk:"old_value"`
}
