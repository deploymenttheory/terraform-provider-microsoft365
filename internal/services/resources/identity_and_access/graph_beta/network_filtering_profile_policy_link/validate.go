package graphBetaNetworkFilteringProfilePolicyLink

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *NetworkFilteringProfilePolicyLinkResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		filteringPolicyOnlyFieldsValidator{},
	}
}

type filteringPolicyOnlyFieldsValidator struct{}

func (v filteringPolicyOnlyFieldsValidator) Description(_ context.Context) string {
	return "priority and logging_state must only be set when policy_type is filtering_policy"
}

func (v filteringPolicyOnlyFieldsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v filteringPolicyOnlyFieldsValidator) ValidateResource(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	if req.Config.Raw.IsNull() {
		return
	}

	var policyType types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("policy_type"), &policyType)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var priority types.Int64
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("priority"), &priority)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var loggingState types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("logging_state"), &loggingState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(validateFilteringPolicyOnlyFields(policyType, priority, loggingState)...)
}

func validateFilteringPolicyOnlyFields(policyType types.String, priority types.Int64, loggingState types.String) diag.Diagnostics {
	var diags diag.Diagnostics

	if policyType.IsNull() || policyType.IsUnknown() || policyType.ValueString() == policyTypeFiltering {
		return diags
	}

	if !priority.IsNull() && !priority.IsUnknown() {
		diags.AddAttributeError(
			path.Root("priority"),
			"priority is only supported for filtering_policy",
			fmt.Sprintf("priority must be omitted when policy_type is %q. Microsoft Graph only accepts priority on filteringPolicyLink requests.", policyType.ValueString()),
		)
	}

	if !loggingState.IsNull() && !loggingState.IsUnknown() {
		diags.AddAttributeError(
			path.Root("logging_state"),
			"logging_state is only supported for filtering_policy",
			fmt.Sprintf("logging_state must be omitted when policy_type is %q. Microsoft Graph only accepts loggingState on filteringPolicyLink requests.", policyType.ValueString()),
		)
	}

	return diags
}
