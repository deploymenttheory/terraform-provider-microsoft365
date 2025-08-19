package graphBetaMacosDeviceConfigurationTemplates

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan validates that exactly one configuration type is specified
func (r *MacosDeviceConfigurationTemplatesResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan MacosDeviceConfigurationTemplatesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Count how many configuration types are specified
	configCount := 0
	if !plan.CustomConfiguration.IsNull() && !plan.CustomConfiguration.IsUnknown() {
		configCount++
	}
	if !plan.PreferenceFile.IsNull() && !plan.PreferenceFile.IsUnknown() {
		configCount++
	}
	if !plan.TrustedCertificate.IsNull() && !plan.TrustedCertificate.IsUnknown() {
		configCount++
	}
	if !plan.ScepCertificate.IsNull() && !plan.ScepCertificate.IsUnknown() {
		configCount++
	}
	if !plan.PkcsCertificate.IsNull() && !plan.PkcsCertificate.IsUnknown() {
		configCount++
	}

	if configCount == 0 {
		resp.Diagnostics.AddError(
			"Missing Configuration",
			"Exactly one of custom_configuration, preference_file, trusted_certificate, scep_certificate, or pkcs_certificate must be specified.",
		)
	} else if configCount > 1 {
		resp.Diagnostics.AddError(
			"Multiple Configurations",
			"Only one of custom_configuration, preference_file, trusted_certificate, scep_certificate, or pkcs_certificate may be specified.",
		)
	}
}
