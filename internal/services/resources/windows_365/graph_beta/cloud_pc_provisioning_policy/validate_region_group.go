package graphBetaCloudPcProvisioningPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// legacyMiddleEastRegionGroupValidator rejects region_group = "middleEast". Windows 365 maps the legacy Middle East
// region group to the Israel Central group (https://learn.microsoft.com/en-us/windows-365/enterprise/enhanced-resiliency-mhn,
// "Updating grouping in existing provisioning policies"), and Graph returns regionGroup "israel" after a create with
// middleEast (#4904), so the configured value can never round-trip.
type legacyMiddleEastRegionGroupValidator struct{}

const legacyMiddleEastRegionGroupDetail = "region_group = \"middleEast\" is the legacy Middle East region group. Windows 365 maps it to the israel (Israel Central) region group, " +
	"so the policy would read back as israel rather than the whole Middle East geography. To select every region in the Middle East geography, set " +
	"geographic_location_type = \"middleEast\", region_group = \"automatic\" and region_name = \"automatic\". To select only Israel Central or UAE, " +
	"set region_group = \"israel\" or region_group = \"uae\"."

func (legacyMiddleEastRegionGroupValidator) Description(context.Context) string {
	return "region_group middleEast is mapped by Windows 365 to israel; select the Middle East geography with geographic_location_type = middleEast and region_group = automatic."
}

func (v legacyMiddleEastRegionGroupValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (legacyMiddleEastRegionGroupValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || req.ConfigValue.ValueString() != "middleEast" {
		return
	}
	resp.Diagnostics.AddAttributeError(req.Path, "Unsupported region group", legacyMiddleEastRegionGroupDetail)
}
