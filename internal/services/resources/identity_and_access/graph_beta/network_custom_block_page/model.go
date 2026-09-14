package graphBetaNetworkCustomBlockPage

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var configurationAttrTypes = map[string]attr.Type{
	"body": types.StringType,
}

type NetworkCustomBlockPageResourceModel struct {
	ID            types.String   `tfsdk:"id"`
	State         types.String   `tfsdk:"state"`
	Configuration types.Object   `tfsdk:"configuration"`
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
}
