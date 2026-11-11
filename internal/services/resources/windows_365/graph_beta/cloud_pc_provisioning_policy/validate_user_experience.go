package graphBetaCloudPcProvisioningPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type userExperienceProvisioningTypeValidator struct{}

func (userExperienceProvisioningTypeValidator) Description(context.Context) string {
	return "cloudApp requires provisioning_type = sharedByEntraGroup."
}
func (v userExperienceProvisioningTypeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}
func (userExperienceProvisioningTypeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.ValueString() != "cloudApp" {
		return
	}
	var provisioningType types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("provisioning_type"), &provisioningType)...)
	if resp.Diagnostics.HasError() || provisioningType.IsUnknown() {
		return
	}
	if provisioningType.ValueString() != "sharedByEntraGroup" {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid provisioning type", "cloudApp requires provisioning_type = sharedByEntraGroup.")
	}
}
