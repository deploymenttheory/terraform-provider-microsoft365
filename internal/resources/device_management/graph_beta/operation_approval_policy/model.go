// resource REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-rbac-operationapprovalpolicy?view=graph-rest-beta
package graphBetaOperationApprovalPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OperationApprovalPolicyResourceModel struct {
	ID                   types.String                            `tfsdk:"id"`
	DisplayName          types.String                            `tfsdk:"display_name"`
	Description          types.String                            `tfsdk:"description"`
	LastModifiedDateTime types.String                            `tfsdk:"last_modified_date_time"`
	PolicyType           types.String                            `tfsdk:"policy_type"`
	PolicyPlatform       types.String                            `tfsdk:"policy_platform"`
	PolicySet            OperationApprovalPolicySetResourceModel `tfsdk:"policy_set"`
	ApproverGroupIds     types.Set                               `tfsdk:"approver_group_ids"`
	Timeouts             timeouts.Value                          `tfsdk:"timeouts"`
}

// PolicySet models
type OperationApprovalPolicySetResourceModel struct {
	PolicyType     types.String `tfsdk:"policy_type"`
	PolicyPlatform types.String `tfsdk:"policy_platform"`
}
