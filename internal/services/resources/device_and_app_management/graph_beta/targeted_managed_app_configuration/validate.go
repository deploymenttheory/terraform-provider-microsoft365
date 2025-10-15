package graphBetaTargetedManagedAppConfigurations

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// validateRequest runs a sequence of validators against the planned configuration
// and returns any diagnostics. Caller is responsible for appending these to the
// response and short-circuiting on error.
func validateRequest(ctx context.Context, plan *TargetedManagedAppConfigurationResourceModel) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	validators := []planValidator{
		newAppGroupTypeValidator(),
	}

	for _, v := range validators {
		tflog.Debug(ctx, "Validate: running validator - "+v.Name())
		if d := v.Validate(ctx, plan); d.HasError() {
			diagnostics.Append(d...)
		}
	}

	return diagnostics
}

// planValidator defines the contract for request validators
type planValidator interface {
	Name() string
	Validate(ctx context.Context, plan *TargetedManagedAppConfigurationResourceModel) diag.Diagnostics
}

// appGroupTypeValidator enforces constraints based on app_group_type
type appGroupTypeValidator struct{}

func newAppGroupTypeValidator() planValidator { return &appGroupTypeValidator{} }

func (v *appGroupTypeValidator) Name() string { return "validateAppGroupType" }

func (v *appGroupTypeValidator) Validate(ctx context.Context, plan *TargetedManagedAppConfigurationResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// If app_group_type is "allApps" then the apps field must NOT be present in HCL
	if !plan.AppGroupType.IsNull() && !plan.AppGroupType.IsUnknown() && plan.AppGroupType.ValueString() == "allApps" {
		// Presence means explicitly set in HCL (even if empty). When omitted, it's typically Null/Unknown.
		if !plan.Apps.IsNull() && !plan.Apps.IsUnknown() {
			diags.AddError(
				"Invalid Configuration",
				"When app_group_type is \"allApps\", the \"apps\" attribute must not be specified in the configuration.",
			)
			return diags
		}
	}

	// Additional checks for other app_group_type values could be added here
	_ = types.StringType // keep types import referenced for future validators

	return diags
}
