package graphBetaCloudPcProvisioningPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type persistenceProvisioningTypeValidator struct{}

func (persistenceProvisioningTypeValidator) Description(context.Context) string {
	return "User settings persistence requires provisioning_type = sharedByEntraGroup."
}
func (v persistenceProvisioningTypeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}
func (persistenceProvisioningTypeValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	var provisioningType types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("provisioning_type"), &provisioningType)...)
	if resp.Diagnostics.HasError() || provisioningType.IsUnknown() {
		return
	}
	if provisioningType.ValueString() != "sharedByEntraGroup" {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid provisioning type", "User settings persistence requires provisioning_type = sharedByEntraGroup.")
	}
}
